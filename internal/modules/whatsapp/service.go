package whatsapp

type Service struct {
	flow    *Flow
	session *SessionService
	client  *Client
}

func NewService(flow *Flow, session *SessionService, client *Client) *Service {
	return &Service{
		flow:    flow,
		session: session,
		client:  client,
	}
}

func (s *Service) HandleIncomingMessage(phone string, message string) {

	session := s.session.GetOrCreate(phone)

	// 2. Process message through flow
	reply := s.flow.Handle(session, message)

	// 3. Save session state
	s.session.Save(session)

	// 4. Send message back to WhatsApp
	_ = s.client.SendText(phone, reply)
	
}
