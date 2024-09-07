package cronservice

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type (
	CronService struct {
		cfg  *configs.Config
		rodr orderrepo.IOrderRepository
	}

	ICronService interface {
		Init() (gocron.Scheduler, error)
	}
)

func NewCronService(
	cfg *configs.Config,
	rodr orderrepo.IOrderRepository,
) *CronService {
	return &CronService{
		cfg:  cfg,
		rodr: rodr,
	}
}

func (s *CronService) Init() (gocron.Scheduler, error) {
	// * get time location from env
	tz, err := time.LoadLocation(s.cfg.API.Timezone)
	if err != nil {
		return nil, err
	}

	// * define a new scheduler instance
	gsch, err := gocron.NewScheduler(
		gocron.WithLocation(tz),
	)

	if err != nil {
		return nil, err
	}

	// * add a order job
	s.orderSchedule(gsch)

	// * start the scheduler
	gsch.Start()

	return gsch, nil
}

func (s *CronService) orderSchedule(gsch gocron.Scheduler) {
	var (
		errs []error
	)

	err := s.scheduleOrderCancelled(gsch)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		utlogger.Error(err)
	}
}

func (s *CronService) scheduleOrderCancelled(gsch gocron.Scheduler) error {
	_, err := gsch.NewJob(
		gocron.DurationJob(
			// TODO: use the env variable
			time.Duration(1)*time.Minute,
		),
		gocron.NewTask(
			func() error {
				err := s.rodr.UpdateAutomaticallyStatus(consttypes.OS_CANCELLED, s.cfg.OrderBuffer.AutomaticallyCancelled, []consttypes.OrderStatus{consttypes.OS_PLACED})
				if err != nil {
					utlogger.Error(err)
					return err
				}

				return nil
			},
		),
	)

	if err != nil {
		utlogger.Error(err)
		return err
	}

	utlogger.Info(fmt.Sprintf("Service for Cron %s Running!", "Update Order Expired"))

	return nil
}
