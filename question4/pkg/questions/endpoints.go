package questions

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/klog/v2"
)

// LoginRequest structure of a login request message
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse structure of a login response message
type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// ProfileData structure of a ProfileData message
type ProfileData struct {
	ProfileID               string `json:"profileId"`
	Name                    string `json:"name"`
	Age                     int    `json:"age"`
	FavoriteColor           string `json:"favoriteColor"`
	FavoriteOperatingSystem string `json:"favoriteOperatingSystem"`
}

// ProfileUpdateRequest structure of a profile update request message
type ProfileUpdateRequest struct {
	Profile *ProfileData `json:"profile,omitempty"`
	Token   string       `json:"token"`
}

// ProfileUpdateResponse structure of a profile update response message
type ProfileUpdateResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// ProfileGetRequest structure of a profile get request message
type ProfileGetRequest struct {
	ProfileID string `json:"profileId"`
	Token     string `json:"token"`
}

// ProfileGetResponse structure of a profile get response message
type ProfileGetResponse struct {
	Profile *ProfileData `json:"profile,omitempty"`
	Status  string       `json:"status"`
	Token   string       `json:"token"`
}

// ProfileDeleteRequest structure of a profile delete request message
type ProfileDeleteRequest struct {
	ProfileID string `json:"profileId"`
	Token     string `json:"token"`
}

// ProfileDeleteResponse structure of a profile delete response message
type ProfileDeleteResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// DefaultPath endpoint to the default path
func (p *Profile) DefaultPath(w http.ResponseWriter, r *http.Request) {
	klog.InfoS("profile.DefaultPath", "EndPoint:", r.URL.Path)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		klog.ErrorS(err, "profile.DefaultPath, ioutil.ReadAll", "URL", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	klog.InfoS("profile.DefaultPath", "body", string(body))

	w.WriteHeader(http.StatusNotFound)
}

// Login endpoint to login
func (p *Profile) Login(w http.ResponseWriter, r *http.Request) {
	klog.V(3).Info("Login")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		klog.ErrorS(err, "Reading body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := &LoginRequest{}
	if err := json.Unmarshal(body, &login); err != nil {
		klog.ErrorS(err, "Error Unmashaling", "body", string(body))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbHashedPassword, ok := p.getUser(login.Username)
	if !ok {
		klog.InfoS("Unknown user", "login", login)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := bcrypt.CompareHashAndPassword(dbHashedPassword, []byte(login.Password)); err != nil {
		klog.ErrorS(err, "Unknown user")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token := guuid.New().String()
	if !p.setSession(token, login.Username) {
		klog.ErrorS(err, "Set Session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := LoginResponse{Token: token, Status: "Success"}
	payload, err := json.Marshal(response)
	if err != nil {
		klog.ErrorS(err, "Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		klog.ErrorS(err, "Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ProfileUpdate endpoint to update a profile
func (p *Profile) ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	klog.V(3).Info("ProfileUpdate")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		klog.ErrorS(err, "Reading body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := &ProfileUpdateRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		klog.ErrorS(err, "Error Unmashaling", "body", string(body))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := p.getSession(request.Token)
	if !ok {
		klog.ErrorS(err, "Session not found", "token", request.Token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if request.Profile == nil {
		klog.ErrorS(err, "Profile missing", "request", request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile, err := json.Marshal(request.Profile)
	if err != nil {
		klog.ErrorS(err, "Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !p.setProfile(user, string(profile)) {
		klog.ErrorS(err, "Set Profile", "user", user)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := ProfileUpdateResponse{Token: request.Token, Status: "Success"}
	payload, err := json.Marshal(response)
	if err != nil {
		klog.ErrorS(err, "Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		klog.ErrorS(err, "Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ProfileDelete endpoint to delete a profile
func (p *Profile) ProfileDelete(w http.ResponseWriter, r *http.Request) {
	klog.V(3).Info("ProfileDelete")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		klog.ErrorS(err, "Reading body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := &ProfileUpdateRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		klog.ErrorS(err, "Error Unmashaling", "body", string(body))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := p.getSession(request.Token)
	if !ok {
		klog.ErrorS(err, "Session not found", "token", request.Token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	response := ProfileDeleteResponse{Token: request.Token, Status: "Success"}
	if _, ok := p.getProfile(user); !ok {
		response.Status = "Failed"
	} else {
		if !p.deleteProfile(user) {
			response.Status = "Failed"
		}
	}

	payload, err := json.Marshal(response)
	if err != nil {
		klog.ErrorS(err, "Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		klog.ErrorS(err, "Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetProfile endpoint to get a profile
func (p *Profile) GetProfile(w http.ResponseWriter, r *http.Request) {
	klog.V(3).Info("GetProfile")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		klog.ErrorS(err, "Reading body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := &ProfileGetRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		klog.ErrorS(err, "Error Unmashaling", "body", string(body))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := p.getSession(request.Token)
	if !ok {
		klog.ErrorS(err, "Session not found", "token", request.Token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	response := ProfileGetResponse{Profile: nil, Token: request.Token, Status: "NotFound"}

	if profile, ok := p.getProfile(user); ok {
		profileData := &ProfileData{}
		if err := json.Unmarshal([]byte(profile), profileData); err != nil {
			klog.ErrorS(err, "Error Unmashaling", "profile", profile)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.Profile = profileData
		response.Status = "Success"
	}

	payload, err := json.Marshal(response)
	if err != nil {
		klog.ErrorS(err, "Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		klog.ErrorS(err, "Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
