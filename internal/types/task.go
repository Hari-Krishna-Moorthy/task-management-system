package types

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255" validateMsg:"Title is required and must be between 1 and 255 characters"`
	Description string `json:"description" validate:"required,min=1,max=512" validateMsg:"Description is required and must be at least 1 character"`
	DueDate     string `json:"due_date" validate:"required,datetime=2006-01-02" validateMsg:"DueDate is required and must be a valid date in YYYY-MM-DD format"`
}

type UpdateTaskRequest struct {
	ID          string `json:"id" validate:"required", validateMsg:"ID is required"`
	Title       string `json:"title" validate:"omitempty,min=1,max=255" validateMsg:"Title is required and must be between 1 and 255 characters"`
	Description string `json:"description" validate:"omitempty,min=1,max=512" validateMsg:"Description is required and must be at least 1 character"`
	DueDate     string `json:"due_date" validate:"omitempty,datetime=2006-01-02" validateMsg:"DueDate is required and must be a valid date in YYYY-MM-DD format"`
	Status      string `json:"status" validate:"omitempty,oneof=ToDo InProgress Done" validateMsg:"Status is required and must be one of 'ToDo', 'InProgress', or 'Done'"`
}

type GetTaskRequest struct {
	ID string `json:"id" validate:"required", validateMsg:"ID is required"`
}

type DeleteTaskRequest struct {
	ID string `json:"id" validate:"required", validateMsg:"ID is required"`
}

type ListTasksRequest struct {
	Status string `json:"status" validate:"omitempty,oneof=ToDo InProgress Done" validateMsg:"Status is required and must be one of 'ToDo', 'InProgress', or 'Done'"`
}
