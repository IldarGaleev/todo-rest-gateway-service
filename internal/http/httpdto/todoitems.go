package httpdto

type TaskItem struct {
	ID     uint64 `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
} //@name TaskItem

type TaskItemChanges struct {
	Title  *string `json:"title,omitempty"`
	IsDone *bool   `json:"is_done,omitempty"`
} //@name TaskItemChanges

type GetTaskListResponse struct {
	GeneralResponse
	Tasks []TaskItem `json:"tasks"`
} //@name GetTaskListResponse

type GetTaskByIDResponse struct {
	GeneralResponse
	Task TaskItem `json:"task"`
} //@name GetTaskByIDResponse
