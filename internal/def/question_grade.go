package def

type GradeName string

const (
	GradeJunior GradeName = "junior"
	GradeMiddle GradeName = "middle"
	GradeSenior GradeName = "senior"
)

func (g GradeName) String() string {
	return string(g)
}

func ValidateGradeName(value string) (GradeName, error) {
	grade := GradeName(value)
	switch grade {
	case GradeJunior, GradeMiddle, GradeSenior:
		return grade, nil
	default:
		return "", ErrInvalidGradeValue
	}
}
