package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message đại diện cho một tin nhắn trong một cuộc trò chuyện.
type Message struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`                        // ID duy nhất của tin nhắn.
	ConversationID uuid.UUID      `gorm:"type:uuid;not null"`                                                      // ID của cuộc trò chuyện mà tin nhắn thuộc về.
	Conversation   Conversation   `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Tham chiếu đến model Conversation.
	SenderID       uuid.UUID      `gorm:"type:uuid;not null"`                                                      // ID của người gửi tin nhắn.
	Sender         User           `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`       // Tham chiếu đến model User của người gửi.
	Content        string         `gorm:"type:text"`                                                               // Nội dung văn bản của tin nhắn.
	CreatedAt      time.Time      `gorm:"autoCreateTime"`                                                          // Thời điểm tạo tin nhắn.
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`                                                          // Thời điểm cập nhật tin nhắn.
	DeletedAt      gorm.DeletedAt `gorm:"index"`                                                                   // Thời điểm xóa tin nhắn (soft delete).
	Status         bool           `gorm:"default:true"`                                                            // trạng thái hiển thị của tin nhắn true : show, false: hide.
	Likes          []LikeMessage  `gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`      // Danh sách các lượt thích tin nhắn.
	Media          []MessageMedia `gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`      // Danh sách media (hình ảnh, video) được đính kèm.
}

// Conversation đại diện cho một cuộc trò chuyện giữa hai người hoặc một nhóm.
type Conversation struct {
	ID        uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`                        // ID duy nhất của cuộc trò chuyện.
	Name      *string              `gorm:"type:varchar(255);default:null"`                                          // Tên của cuộc trò chuyện nhóm (null nếu là tin nhắn 1:1).
	IsGroup   bool                 `gorm:"default:false"`                                                           // True nếu là nhóm, false nếu là tin nhắn 1:1.
	CreatedAt time.Time            `gorm:"autoCreateTime"`                                                          // Thời điểm tạo cuộc trò chuyện.
	UpdatedAt time.Time            `gorm:"autoUpdateTime"`                                                          // Thời điểm cập nhật cuộc trò chuyện.
	DeletedAt gorm.DeletedAt       `gorm:"index"`                                                                   // Thời điểm xóa cuộc trò chuyện (soft delete).
	Members   []ConversationMember `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Danh sách thành viên trong cuộc trò chuyện.
	Messages  []Message            `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Danh sách tin nhắn trong cuộc trò chuyện.
	AvatarUrl string               `gorm:"type:varchar(255);default:null"`                                          // Đường dẫn avatar của nhóm.
}

// ConversationMember đại diện cho một thành viên trong cuộc trò chuyện.
type ConversationMember struct {
	ConversationID uuid.UUID    `gorm:"type:uuid;primary_key;not null"`                                          // ID của cuộc trò chuyện.
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Tham chiếu đến model Conversation.
	UserID         uuid.UUID    `gorm:"type:uuid;primary_key;not null"`                                          // ID của người dùng.
	User           User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`         // Tham chiếu đến model User.
}

// LikeMessage đại diện cho một lượt thích trên một tin nhắn.
type LikeMessage struct {
	UserID    uuid.UUID `gorm:"type:uuid;primary_key;not null"`                                     // ID của người dùng thích tin nhắn.
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`    // Tham chiếu đến model User.
	MessageID uuid.UUID `gorm:"type:uuid;primary_key;not null"`                                     // ID của tin nhắn được thích.
	Message   Message   `gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Tham chiếu đến model Message.
}

// MessageMedia đại diện cho media (hình ảnh, video) được đính kèm vào tin nhắn.
type MessageMedia struct {
	ID        uint      `gorm:"type:int;auto_increment;primary_key"` // ID duy nhất của media.
	MessageID uuid.UUID `gorm:"type:uuid;not null"`                  // ID của tin nhắn mà media được đính kèm.
	MediaUrl  string    `gorm:"type:varchar(255);not null"`          // Đường dẫn đến file media.
	Status    bool      `gorm:"default:true"`                        // trạng thái hiển thị của media true: show, false: hide.
	CreatedAt time.Time `gorm:"autoCreateTime"`                      // Thời gian tạo.
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                      // thời gian cập nhật.
}
