package withdrawal

import "github.com/stylebookgo/agbadevemma/internal/utils"

type WithdrawalWriter interface {
	HandleTransferEvent(event utils.PaystackWebhookEvent) error
}
