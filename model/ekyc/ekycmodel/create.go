package ekycmodel

import (
	"mime/multipart"
)

type CreateKYCProfileRequest struct {
	FrontDocument *multipart.FileHeader `form:"front_document"`
	BackDocument  *multipart.FileHeader `form:"back_document"`
	FaceImage     *multipart.FileHeader `form:"face_image"`
	UserID        uint32                `form:"user_id"`
	CardID        string                `form:"card_id"`
	FullName      string                `form:"full_name"`
	DOB           string                `form:"dob"`
	Gender        string                `form:"gender"`
	Address       string                `form:"address"`
	ClientSession string                `form:"client_session"`
}
