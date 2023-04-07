package auth

import "testing"

func TestAuthAPI(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		setupMocks    func()
		checkResponse func()
	}
}

func randomUser(t *testing.T) (*user, error) {

}
