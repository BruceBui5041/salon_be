package appconst

// DynamoDB
const (
	VideoURLPrefix   = "video_url"
	UserPrefix       = "user"
	EnrollmentPrefix = "enroll"
)

// local TOPIC
const (
	TopicNewVideoUploaded           = "new_video_uploaded"
	TopicUpdateVideoProcessingState = "update_video_processing_state"
	TopicVideoProcessed             = "video_processed"
	TopicEnrollmentChange           = "enrollment_change"
	TopicUserUpdated                = "user_updated"
	TopicReceivedWSMsg              = "revevied_ws_msg"
	TopicBookingEvent               = "booking_event"
)

// sqs message group id
const (
	ReqProcessVideoGroupId            = "req-process-video"
	UpdateVideoProcessProgressGroupId = "update-video-process-progress"
)

// token
const (
	AccessTokenName = "access_token"
	TokenExpiry     = 60 * 60 * 24 * 7
)

const (
	MaxConcurrentS3Push  = 50
	AWSVideoS3BuckerName = "hls-video-segment"
	AWSCloudFrontVideo   = "https://d17cfikyg12m49.cloudfront.net"

	AWSPublicBucket         = "salon-dev-pub"
	AWSCloudFrontPublicFile = "https://d3i048dqjftjb3.cloudfront.net"
)
