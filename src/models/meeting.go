package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/MakotoNakai/lets-schedule/config"
)

type Meeting struct {
	Id int `gorm:"primaryKey:not null:autoIncrement" json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	IsOnsite bool `json:"is_onsite"`
	IsOnline bool `json:"is_online"`
	Place string `json:"place"`
	Url string `json:"url"`
	AllParticipantsResponded bool `json:"all_participants_responded"`
	IsConfirmed bool `json:"is_confirmed"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
	ActualStartTime time.Time `json:"actual_start_time"`
	ActualEndTime time.Time `json:"actual_end_time"`
	Hour int `json:"hour"`
	CreatedAt time.Time `gorm:"autoCreateTime:int" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:int" json:"updated_at"`
}

func IsTitleEmpty(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	return m_.Title == "", nil
}

func IsHourEmpty(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	if m_.Hour < 0 {
		return false, config.ErrMeetingHourIsNegative
	}
	return m_.Hour == 0, nil
}

func IsOnsiteButNoPlaceSpecified(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	return m_.IsOnsite == true && m_.IsOnline == false && m_.Place == "", nil
}

func IsOnlineButNoURLSpecified(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	return m_.IsOnsite == false && m_.IsOnline == true && m_.Url == "", nil
}

func IsHybridButNeitherPlaceOrURLSpecified(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	return m_.IsOnsite == true && m_.IsOnline == true && m_.Place == "" && m_.Url == "", nil
}

func IsHybridButNoPlaceSpecified(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	return m_.IsOnsite == true && m_.IsOnline == true && m_.Place == "", nil
}

func IsHybridButNoURLSpecified(m *Meeting) (bool, error) {
	if m == nil {
		return false, config.ErrMeetingNotFound
	}
	m_ := *m
	return m_.IsOnsite == true && m_.IsOnline == true && m_.Url == "", nil
}

func GetMeetingById(db *gorm.DB, Id int) (Meeting, error) {
	meeting := Meeting{}
	err := db.Table("meetings").
		Select("*").
		Where("meetings.id = ?", Id).
		Find(&meeting).Error
	if err != nil {
		return meeting, config.ErrRecordNotFound
	}
	return meeting, nil
}

func GetMeetingsByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("users.id = ?", UserId).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}

func GetConfirmedMeetingsForHostByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("participants.user_id = ? AND participants.is_host = ?", UserId, 1).
		Where("meetings.is_confirmed = ?", 1).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}

func GetNotConfirmedMeetingsForHostByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("participants.user_id = ? AND participants.is_host = ?", UserId, 1).
		Where("meetings.is_confirmed = ? AND meetings.all_participants_responded = ?", 0, 1).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}

func GetNotRespondedMeetingsForHostByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("participants.user_id = ? AND participants.is_host = ?", UserId, 1).
		Where("meetings.is_confirmed = ? AND meetings.all_participants_responded = ?", 0, 0).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}

func GetConfirmedMeetingsForGuestByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("participants.user_id = ? AND participants.is_host = ?", UserId, 0).
		Where("meetings.is_confirmed = ?", 1).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}

func GetNotConfirmedMeetingsForGuestByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("participants.user_id = ? AND participants.is_host = ? AND participants.has_responded = ?", UserId, 0, 1).
		Where("meetings.is_confirmed = ?", 0).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}

func GetNotRespondedMeetingsForGuestByUserId(db *gorm.DB, UserId int) ([]Meeting, error) {
	meetings := []Meeting{}
	err := db.Table("meetings").
		Select("*").
		Joins("inner join participants on participants.meeting_id = meetings.id").
		Joins("inner join users on users.id = participants.user_id").
		Where("participants.user_id = ? AND participants.is_host = ? AND participants.has_responded = ?", UserId, 0, 0).
		Where("meetings.is_confirmed = ?", 0).
		Find(&meetings).Error
	if err != nil {
		return meetings, config.ErrRecordNotFound
	}
	return meetings, nil
}


