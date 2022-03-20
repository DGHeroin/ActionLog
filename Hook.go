package ActionLog

type (
    Hook interface {
        Fire(*Entry) error
    }
)
