package models

import (
	"errors"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/enums"
)

type Task struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`

	Completed bool             `json:"completed" bson:"completed"`
	Status    enums.TaskStatus `json:"status" bson:"status"`
	DueDate   time.Time        `json:"due_date" bson:"due_date"`

	UserID string `json:"user_id" bson:"user_id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at" `
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" `
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at" `
}

func (task *Task) GetCollectionsName() string {
	return "tasks"
}

var transitionsMap = map[enums.TaskStatus][]enums.TaskStatus{
	enums.ToDo:       {enums.InProgress},
	enums.InProgress: {enums.Done},
	enums.Done:       {enums.ToDo, enums.InProgress},
}

type transitionRules []func(*Task) error

var transitionRulesMap = map[string]transitionRules{
	"InProgress->Done": {closeTask},
	"Done->InProgress": {reopenTask},
	"Done->ToDo":       {reopenTask},
}

func closeTask(task *Task) error {
	task.Completed = true
	return nil
}

func reopenTask(task *Task) error {
	task.Completed = false
	return nil
}

func (task *Task) CanMoveStatus(currentStatus, newStatus enums.TaskStatus) bool {
	allowedStatuses, exists := transitionsMap[currentStatus]
	if !exists {
		return false
	}
	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}
	return false
}

func (task *Task) ApplyStatusChange(newStatus enums.TaskStatus) error {
	if !task.CanMoveStatus(task.Status, newStatus) {
		return errors.New("invalid status transition from " + task.Status.String() + " to " + newStatus.String())
	}

	ruleKey := task.Status.String() + "->" + newStatus.String()
	if rule, exists := transitionRulesMap[ruleKey]; exists {
		for _, rule := range rule {
			if err := rule(task); err != nil {
				return err
			}
		}
	}

	task.Status = newStatus
	return nil
}
