package job

import (
	"log"
	"time"

	"piemdm/internal/service"
)

type Scanner struct {
	// lastID      uint64
	userService service.UserService
	cronService service.CronService
}

const (
	ScannerSize = 10
)

func NewScanner(userService service.UserService, cronService service.CronService) *Scanner {
	return &Scanner{
		userService: userService,
		cronService: cronService,
	}
}

func (s *Scanner) Run() {
	err := s.scannerDB()
	if err != nil {
		log.Println("scannerDB", "err", err)
		return
	}
}

func (s *Scanner) scannerDB() error {
	log.Println("Cron Job", "time", time.Now())
	user, err := s.userService.Get(223)
	// where := map[string]any{}
	// where["status"] = "Normal"
	// jobs, err := s.cronService.Find("", where)
	log.Println("user", "user", user, "err", err)
	// s.reset()
	// flag := false
	// for {
	// 	users, err := s.user.MGet(s.lastID, ScannerSize)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(users) < ScannerSize {
	// 		flag = true
	// 	}
	// 	s.lastID = users[len(users)-1].ID
	// 	for k, v := range users {
	// 		logger.Info("k, v", "k", k, "v", v)
	// 	}
	// 	if flag {
	// 		return nil
	// 	}
	// }
	return nil
}

// func (s *Scanner) reset() {
// 	s.lastID = 0
// }
