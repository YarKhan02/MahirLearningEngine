package quiz

type CreateQuizRequest struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description,omitempty"`
	Questions   []CreateQuestionRequest `json:"questions"`
}

type CreateQuestionRequest struct {
	Prompt        string                `json:"prompt"`
	Type          string                `json:"type"` // "mcq" | "typed"
	Marks         int                   `json:"marks"`
	AllowMultiple bool                  `json:"allowMultiple,omitempty"`
	Options       []CreateOptionRequest `json:"options,omitempty"`
}

type CreateOptionRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type OptionResponse struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	OrderNo   int    `json:"orderNo"`
	IsCorrect *bool  `json:"isCorrect,omitempty"` // omitted for students
}

type QuestionResponse struct {
	ID            string           `json:"id"`
	Prompt        string           `json:"prompt"`
	Type          string           `json:"type"`
	Marks         int              `json:"marks"`
	AllowMultiple bool             `json:"allowMultiple"`
	OrderNo       int              `json:"orderNo"`
	Options       []OptionResponse `json:"options"`
}

type QuizResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	TotalMarks  int                `json:"totalMarks"`
	CreatedAt   string             `json:"createdAt"`
	Questions   []QuestionResponse `json:"questions"`
}

type AdminQuizResponse struct {
	QuizResponse
	SubmissionCount int `json:"submissionCount"`
	PendingCount    int `json:"pendingCount"`
}

type StudentQuizResponse struct {
	QuizResponse
	Submission *StudentSubmissionResponse `json:"submission,omitempty"`
}

type StudentSubmissionResponse struct {
	ID          string                  `json:"id"`
	Status      string                  `json:"status"`
	Score       int                     `json:"score"`
	TotalMarks  int                     `json:"totalMarks"`
	Remarks     *string                 `json:"remarks,omitempty"`
	SubmittedAt string                  `json:"submittedAt"`
	Answers     []StudentAnswerResponse `json:"answers"`
}

type StudentAnswerResponse struct {
	QuestionID        string   `json:"questionId"`
	AnswerText        string   `json:"answerText,omitempty"`
	SelectedOptionIds []string `json:"selectedOptionIds,omitempty"`
	AwardedMarks      *int     `json:"awardedMarks,omitempty"`
}

type SubmitQuizRequest struct {
	Answers []SubmitAnswerRequest `json:"answers"`
}

type SubmitAnswerRequest struct {
	QuestionID        string   `json:"questionId"`
	AnswerText        string   `json:"answerText,omitempty"`
	SelectedOptionIds []string `json:"selectedOptionIds,omitempty"`
}

/* ---------------- Submissions review + grade (admin) ---------------- */

type SubmissionRowResponse struct {
	ID           string `json:"id"`
	StudentID    string `json:"studentId"`
	StudentName  string `json:"studentName"`
	StudentEmail string `json:"studentEmail"`
	Status       string `json:"status"`
	Score        int    `json:"score"`
	SubmittedAt  string `json:"submittedAt"`
}

type SubmissionSummaryResponse struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Graded  int `json:"graded"`
}

// SubmissionDetailResponse is the full grading view for one submission.
type SubmissionDetailResponse struct {
	ID          string                    `json:"id"`
	QuizTitle   string                    `json:"quizTitle"`
	StudentName string                    `json:"studentName"`
	Status      string                    `json:"status"`
	Score       int                       `json:"score"`
	TotalMarks  int                       `json:"totalMarks"`
	Remarks     *string                   `json:"remarks,omitempty"`
	SubmittedAt string                    `json:"submittedAt"`
	Questions   []GradedQuestionResponse  `json:"questions"`
}

type GradedQuestionResponse struct {
	QuestionID        string           `json:"questionId"`
	Prompt            string           `json:"prompt"`
	Type              string           `json:"type"`
	Marks             int              `json:"marks"`
	AnswerText        string           `json:"answerText,omitempty"`
	SelectedOptionIds []string         `json:"selectedOptionIds,omitempty"`
	AwardedMarks      *int             `json:"awardedMarks,omitempty"`
	Options           []OptionResponse `json:"options"` // includes isCorrect for the teacher
}

type GradeQuizRequest struct {
	Remarks string             `json:"remarks,omitempty"`
	Marks   []GradeAnswerInput `json:"marks"`
}

type GradeAnswerInput struct {
	QuestionID string `json:"questionId"`
	Marks      int    `json:"marks"`
}
