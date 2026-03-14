package scheduler

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/saintbyte/home-ctrl/internal/config"
)

type TaskExecutor func(taskName string)

type Scheduler struct {
	cron     *cron.Cron
	config   *config.Config
	executor TaskExecutor
	taskIDs  map[string]cron.EntryID
	mu       sync.Mutex
}

func NewScheduler(cfg *config.Config, executor TaskExecutor) *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		config:   cfg,
		executor: executor,
		taskIDs:  make(map[string]cron.EntryID),
	}
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.config.Tasks {
		if task.Enabled && task.Schedule != "" {
			s.addTask(task)
		}
	}
	s.cron.Start()
	fmt.Printf("Scheduler started with %d tasks\n", len(s.taskIDs))
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.cron.Stop()
	<-ctx.Done()
	fmt.Println("Scheduler stopped")
}

func (s *Scheduler) addTask(task config.Task) {
	id, err := s.cron.AddFunc(task.Schedule, func() {
		fmt.Printf("Running task: %s\n", task.Name)
		if s.executor != nil {
			s.executor(task.Name)
		}
	})

	if err != nil {
		fmt.Printf("Failed to add task %s: %v\n", task.Name, err)
		return
	}

	s.taskIDs[task.Name] = id
	fmt.Printf("Added task: %s with schedule: %s\n", task.Name, task.Schedule)
}

func (s *Scheduler) EnableTask(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.config.Tasks {
		if task.Name == name {
			task.Enabled = true
			s.addTask(task)
			return nil
		}
	}
	return fmt.Errorf("task not found: %s", name)
}

func (s *Scheduler) DisableTask(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, ok := s.taskIDs[name]; ok {
		s.cron.Remove(id)
		delete(s.taskIDs, name)

		for i := range s.config.Tasks {
			if s.config.Tasks[i].Name == name {
				s.config.Tasks[i].Enabled = false
				break
			}
		}
		return nil
	}
	return fmt.Errorf("task not found: %s", name)
}

func (s *Scheduler) RunTask(name string) error {
	for _, task := range s.config.Tasks {
		if task.Name == name {
			go func() {
				fmt.Printf("Manually running task: %s\n", name)
				if s.executor != nil {
					s.executor(name)
				}
			}()
			return nil
		}
	}
	return fmt.Errorf("task not found: %s", name)
}

func (s *Scheduler) GetTasks() []config.Task {
	return s.config.Tasks
}
