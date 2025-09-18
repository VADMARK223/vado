package http

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"vado/gui/common"
	"vado/gui/constant"
	"vado/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const initBalance = 50
const initBank = 0         // Стартовый капитал в банке
const balanceDeltaGui = 30 // Изменение на которое изменяется баланс из GUI// Стартовый капитал клиента

var (
	moneyMtx   sync.Mutex
	balance    atomic.Int64  // Денег у клиента
	bank       atomic.Int64  // Денег в банке
	balanceLbl *widget.Label // Информация о балансе клиента
	bankLbl    *widget.Label // Информация о балансе банка
)

func init() {
	balance.Store(initBalance)
	bank.Store(initBank)
}

func createMoneyGui() fyne.CanvasObject {
	title := canvas.NewText("Управление балансом", constant.Gold())
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	balanceLbl = widget.NewLabel("")
	labelCentered := container.New(layout.NewCenterLayout(), balanceLbl)
	updateBalanceText()

	balanceDecreaseBtn := common.CreateBtn("-", nil, func() {
		decreaseBalance(balanceDeltaGui)
		updateBalanceText()
	})
	btnDecreaseContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(60, 40)), balanceDecreaseBtn)
	balanceDecreaseBtn.Resize(fyne.NewSize(100, 40))
	balanceIncreaseBtn := common.CreateBtn("+", nil, func() {
		balance.Add(balanceDeltaGui)
		updateBalanceText()
	})
	btnIncreaseContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(60, 40)), balanceIncreaseBtn)
	balanceBox := container.NewHBox(labelCentered, layout.NewSpacer(), btnDecreaseContainer, btnIncreaseContainer)

	bankLbl = widget.NewLabel("")
	updateBankText()
	bankBox := container.NewHBox(bankLbl)

	vBox := container.NewVBox(title, balanceBox, bankBox)
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 255, A: 50})
	content := container.NewStack(bg, vBox)
	return content
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		msg := "Error: Method not allowed."
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(msg)
		if _, err := w.Write([]byte(msg)); err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	httpRequestBody, errRequestBody := io.ReadAll(r.Body)
	if errRequestBody != nil {
		fmt.Println("Fail to read HTTP request body:", errRequestBody)
		return
	}

	httpRequestBodyStr := string(httpRequestBody)

	paymentAmount, paymentAmountErr := strconv.Atoi(httpRequestBodyStr)
	if paymentAmountErr != nil {
		msg := "Fail to parse payment amount:" + paymentAmountErr.Error()
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(msg)); err != nil {
			fmt.Println(msg)
		}
		fmt.Println(msg)
		return
	}
	fmt.Println("Payment amount:", paymentAmount)

	moneyMtx.Lock()
	decreaseResult := decreaseBalance(int64(paymentAmount))
	msg := func() string {
		if decreaseResult {
			updateBalanceText()
			return util.Tpl("Successful payment: %s$, current balance: %d$", httpRequestBodyStr, balance.Load())
		}
		return util.Tpl("FAIL: Payment amount out of balance: %d$", paymentAmount)
	}()
	moneyMtx.Unlock()

	fmt.Println("Balance after pay:", balance.Load())

	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Successful payment")
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, errRequestBody := io.ReadAll(r.Body)
	if errRequestBody != nil {
		fmt.Println("Fail to read HTTP request body:", errRequestBody)
		return
	}

	httpRequestBodyStr := string(httpRequestBody)

	saveAmount, saveAmountErr := strconv.Atoi(httpRequestBodyStr)
	if saveAmountErr != nil {
		msg := util.Tpl("Fail to parse save amount: %s", saveAmountErr)
		fmt.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(msg)); err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	fmt.Println("Save amount:", saveAmount)

	moneyMtx.Lock()
	decreaseResult := decreaseBalance(int64(saveAmount))
	msg := func() string {
		if decreaseResult {
			bank.Add(int64(saveAmount))
			updateBalanceText()
			updateBankText()
			return util.Tpl("Success save: %s$, current balance: %d$", httpRequestBodyStr, balance.Load())
		}
		return util.Tpl("FAIL: Save amount out of balance: %d$", saveAmount)
	}()
	moneyMtx.Unlock()

	fmt.Println("Balance after save:", balance.Load())

	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Successful save")
	}
}

func updateBalanceText() {
	if a := fyne.CurrentApp(); a != nil {
		fyne.Do(func() {
			balanceLbl.SetText(util.Tpl("Баланс: %d$", balance.Load()))
		})
	}
}

func updateBankText() string {
	if a := fyne.CurrentApp(); a != nil {
		fyne.Do(func() {
			bankLbl.SetText(util.Tpl("Bank: %d$", bank.Load()))
		})
	}
	return util.Tpl("Банк: %d$", bank.Load())
}

func decreaseBalance(value int64) bool {
	if value <= balance.Load() {
		balance.Add(-value)
		return true
	} else {
		return false
	}
}
