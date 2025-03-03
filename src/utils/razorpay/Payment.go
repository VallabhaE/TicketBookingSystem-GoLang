package razorpay

import (
	"github.com/razorpay/razorpay-go"
	"github.com/razorpay/razorpay-go/utils"
)

func VerifyPayment(orderId, paymentId, signature string) bool {
	params := map[string]any{
		"razorpay_order_id":   orderId,
		"razorpay_payment_id": paymentId,
	}

	// "order_Q2GAynRr6DPsPU"
	// "pay_Q2GFJt0BIUg6Nx"
	// signature := "e310fb0ed23da20e3d4247362700397677bdce1006996ac0a63991666352950b";
	secret := "9eRxQ9TgUzsEDBoo3GpgIXgA"
	return utils.VerifyPaymentSignature(params, signature, secret)
}


func CreateOrderId(amount int,receipt int) (map[string]interface{},error){
	client := razorpay.NewClient("rzp_test_FbANxXzBiqxy1q", "9eRxQ9TgUzsEDBoo3GpgIXgA")
	data := map[string]any{
		"amount":   amount * 100, // Amount is in currency subunits. Default currency is INR. Hence, 50000 refers to 50000 paise
		"currency": "INR",
		"receipt":  receipt,
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return nil,err
	}

	return body,nil
}