package errno

var (
	// ErrEmailAlreadyExist 代表邮箱已经存在.
	ErrEmailAlreadyExist = &Errno{HTTP: 400, Code: 20201, Message: "Email already exist."}

	// ErrUserNotFound 表示未找到用户.
	ErrUserNotFound = &Errno{HTTP: 404, Code: 20202, Message: "User was not found."}

	// ErrPasswordIncorrect 表示密码不正确.
	ErrPasswordIncorrect = &Errno{HTTP: 401, Code: 20203, Message: "Password was incorrect."}
)
