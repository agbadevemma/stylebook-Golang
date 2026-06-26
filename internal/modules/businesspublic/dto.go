package businesspublic

import "github.com/google/uuid"

type PublicBusinessResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Timezone string    `json:"timezone"`
}
