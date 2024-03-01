package profile

type Aggregate struct {
	Ehid             string `json:"ehid"`
	EmployeeId       string `json:"employee_id"`
	Name             string `json:"name"`
	EmailAddress     string `json:"email_address"`
	Dob              string `json:"dob"`
	Grade            string `json:"grade"`
	Title            string `json:"title"`
	OrganizationNode string `json:"organization_node"`
}
