package auth

type Instruction string

const (
	BalanceQuery         Instruction = "balanceQuery"
	DepositAddressQuery  Instruction = "depositAddressQuery"
	DepositQueryAll      Instruction = "depositQueryAll"
	FillHistoryQueryAll  Instruction = "fillHistoryQueryAll"
	OrderCancel          Instruction = "orderCancel"
	OrderCancelAll       Instruction = "orderCancelAll"
	OrderExecute         Instruction = "orderExecute"
	OrderHistoryQueryAll Instruction = "orderHistoryQueryAll"
	OrderQuery           Instruction = "orderQuery"
	OrderQueryAll        Instruction = "orderQueryAll"
	Withdraw             Instruction = "withdraw"
	WithdrawalQueryAll   Instruction = "withdrawalQueryAll"
)
