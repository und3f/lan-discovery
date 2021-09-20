package scanner

import (
	"net"
	"testing"
)

type TestObserver struct {
	ID            int
	Notifications []*Host
}

func (p *TestObserver) Update(h *Host) *Host {
	p.Notifications = append(p.Notifications, h)
	return nil
}

const Host1Str = "192.168.0.1"
const Host2Str = "192.168.0.2"

func TestObservers(t *testing.T) {
	testObserver1 := &TestObserver{1, nil}
	testObserver2 := &TestObserver{2, nil}
	testObserver3 := &TestObserver{3, nil}

	publisher := &HostUpdatePublisher{}

	t.Run("Subscribe()", func(t *testing.T) {
		publisher.Subscribe(testObserver1)
		publisher.Subscribe(testObserver2)
		publisher.Subscribe(testObserver3)

		if want, got := 3, len(publisher.subscribers); want != got {
			t.Errorf("The size of the subscribers list is not the expected. %d != %d\n", want, got)
		}
	})

	t.Run("NotifySubscribers()", func(t *testing.T) {
		host1 := NewHost(net.ParseIP(Host1Str))
		publisher.NotifySubscribers(host1)

		if want, got := 1, len(testObserver1.Notifications); want != got {
			t.Errorf("The amount of received notifications is not the expected. %d != %d\n", want, got)
		}

		if want, got := 1, len(testObserver2.Notifications); want != got {
			t.Errorf("The amount of received notifications is not the expected. %d != %d\n", want, got)
		}

		if want, got := 1, len(testObserver2.Notifications); want != got {
			t.Errorf("The amount of received notifications is not the expected. %d != %d\n", want, got)
		}

		if want, got := Host1Str, testObserver1.Notifications[0].IP.String(); want != got {
			t.Errorf("Wrong notification data. %s != %s\n", want, got)
		}
	})

	t.Run("Unsubscribe()", func(t *testing.T) {
		publisher.Unsubscribe(testObserver2)

		if want, got := 2, len(publisher.subscribers); want != got {
			t.Errorf("The size of the subscribers list is not the expected. %d != %d\n", want, got)
		}

		host2 := NewHost(net.ParseIP(Host2Str))
		publisher.NotifySubscribers(host2)

		if want, got := 2, len(testObserver1.Notifications); want != got {
			t.Errorf("The amount of received notifications is not the expected. %d != %d\n", want, got)
		}
		if want, got := 1, len(testObserver2.Notifications); want != got {
			t.Errorf("The amount of received notifications is not the expected. %d != %d\n", want, got)
		}
		if want, got := 2, len(testObserver3.Notifications); want != got {
			t.Errorf("The amount of received notifications is not the expected. %d != %d\n", want, got)
		}

		if want, got := Host2Str, testObserver1.Notifications[1].IP.String(); want != got {
			t.Errorf("Wrong notification data. %s != %s\n", want, got)
		}
	})
}
