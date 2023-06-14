package handlers

type Transfer struct {
  Id          int
  Origin      string
  Destination string
  Amount      float64
  Reference   string
  Status      string
  TWorkflowId string
  TRunId      string
  TTaskQueue  string
  TInfo       string
}

