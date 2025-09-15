package http

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	c "vado/constant"
	"vado/gui/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	moneyMtx   sync.Mutex
	balance    atomic.Int64  // Денег у клиента
	bank       atomic.Int64  // Денег в банке
	balanceLbl *widget.Label // Информация о балансе клиента
	bankLbl    *widget.Label // Информация о балансе банка
)

func init() {
	balance.Store(c.InitBalance)
	bank.Store(c.InitBank)
}

func createMoneyGui() fyne.CanvasObject {
	title := canvas.NewText("Money", c.Gold())
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	balanceLbl = widget.NewLabel("")
	updateBalanceText()

	balanceDecreaseBtn := common.CreateBtn("-", nil, func() {
		decreaseBalance(c.BalanceDeltaGui)
		updateBalanceText()
	})
	balanceIncreaseBtn := common.CreateBtn("+", nil, func() {
		balance.Add(c.BalanceDeltaGui)
		updateBalanceText()
	})
	balanceBox := container.NewHBox(balanceLbl, balanceDecreaseBtn, balanceIncreaseBtn)

	bankLbl = widget.NewLabel("")
	updateBankText()
	bankBox := container.NewHBox(bankLbl)

	vBox := container.NewVBox(title, balanceBox, bankBox)
	bg := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 255, A: 50})
	content := container.NewStack(bg, vBox)
	return content
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, errRequestBody := io.ReadAll(r.Body)
	if errRequestBody != nil {
		fmt.Println("Fail to read HTTP request body:", errRequestBody)
		return
	}

	httpRequestBodyStr := string(httpRequestBody)

	paymentAmount, paymentAmountErr := strconv.Atoi(httpRequestBodyStr)
	if paymentAmountErr != nil {
		fmt.Println("Fail to parse payment amount:", paymentAmountErr)
		return
	}
	fmt.Println("Payment amount:", paymentAmount)

	moneyMtx.Lock()
	decreaseResult := decreaseBalance(int64(paymentAmount))
	msg := func() string {
		if decreaseResult {
			updateBalanceText()
			return fmt.Sprintf("Successful payment: %s$, current balance: %d$", httpRequestBodyStr, balance.Load())
		}
		return fmt.Sprintf("FAIL: Payment amount out of balance: %d$", paymentAmount)
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
		fmt.Println("Fail to parse save amount:", saveAmountErr)
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
			return fmt.Sprintf("Success save: %s$, current balance: %d$", httpRequestBodyStr, balance.Load())
		}
		return fmt.Sprintf("FAIL: Save amount out of balance: %d$", saveAmount)
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
			balanceLbl.SetText(fmt.Sprintf("Balance: %d$", balance.Load()))
		})
	}
}

func updateBankText() string {
	if a := fyne.CurrentApp(); a != nil {
		fyne.Do(func() {
			bankLbl.SetText(fmt.Sprintf("Bank: %d$", bank.Load()))
		})
	}
	return fmt.Sprintf("Bank: %d$", bank.Load())
}

func decreaseBalance(value int64) bool {
	if value <= balance.Load() {
		balance.Add(-value)
		return true
	} else {
		return false
	}
}
