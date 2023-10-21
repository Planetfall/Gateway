package service

// reportErrorCallback is the callback used by the service controllers
// when an error is encountered in the handlers
func (s *Service) reportErrorCallback(err error) {

	message := "received error from controller"
	s.srv.Raise(message, err, nil)
}
