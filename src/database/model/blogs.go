package model

type Blogs struct {
	BaseModel
	UserId         int64              `gorm:"column:user_id;not null"`
	Title          string             `gorm:"column:title;not null"`
	Content        string             `gorm:"column:content;type:text;not null"`
	IsAnonymous    bool               `gorm:"column:is_anonymous;not null"`
	TotalReviews   int                `gorm:"column:total_reviews;not null"`
	Classification BlogClassification `gorm:"column:classification;not null"`
	Stars          int                `gorm:"column:stars;not null"`
}

func (*Blogs) TableName() string {
	return "blogs"
}

func NewBlogs() *Blogs {
	return &Blogs{
		BaseModel: *NewBaseModel(),
	}
}

func (model *Blogs) CreateData(userId int64, title, content string, isAnonymous bool) Blogs {
	model.UserId = userId
	model.Title = title
	model.Content = content
	model.IsAnonymous = isAnonymous
	return *model
}
