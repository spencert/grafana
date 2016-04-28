package sqlstore

import (
	"testing"

	m "github.com/grafana/grafana/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAlertingStateAccess(t *testing.T) {
	Convey("Test alerting state changes", t, func() {
		InitTestDB(t)

		testDash := insertTestDashboard("dashboard with alerts", 1, "alert")

		items := []m.AlertRule{
			{
				PanelId:     1,
				DashboardId: testDash.Id,
				OrgId:       testDash.OrgId,
				Query:       "Query",
				QueryRefId:  "A",
				WarnLevel:   "> 30",
				CritLevel:   "> 50",
				Interval:    "10",
				Title:       "Alerting title",
				Description: "Alerting description",
				QueryRange:  "5m",
				Aggregator:  "avg",
			},
		}

		cmd := m.SaveAlertsCommand{
			Alerts:      &items,
			DashboardId: testDash.Id,
			OrgId:       1,
			UserId:      1,
		}

		err := SaveAlerts(&cmd)
		So(err, ShouldBeNil)

		Convey("Cannot insert invalid states", func() {
			err = SetNewAlertState(&m.UpdateAlertStateCommand{
				AlertId:  1,
				NewState: "maybe ok",
				Info:     "Shit just hit the fan",
			})

			So(err, ShouldNotBeNil)
		})

		Convey("Changes state to alert", func() {

			err = SetNewAlertState(&m.UpdateAlertStateCommand{
				AlertId:  1,
				NewState: "ALERT",
				Info:     "Shit just hit the fan",
			})

			Convey("can get new state for alert", func() {
				query := &m.GetAlertByIdQuery{Id: 1}
				err := GetAlertById(query)
				So(err, ShouldBeNil)
				So(query.Result.State, ShouldEqual, "ALERT")
			})

			Convey("Changes state to ok", func() {
				err = SetNewAlertState(&m.UpdateAlertStateCommand{
					AlertId:  1,
					NewState: "OK",
					Info:     "Shit just hit the fan",
				})

				Convey("get ok state for alert", func() {
					query := &m.GetAlertByIdQuery{Id: 1}
					err := GetAlertById(query)
					So(err, ShouldBeNil)
					So(query.Result.State, ShouldEqual, "OK")
				})

				Convey("should have two event state logs", func() {
					query := &m.GetAlertsStateLogCommand{
						AlertId: 1,
						OrgId:   1,
					}

					err := GetAlertStateLogByAlertId(query)
					So(err, ShouldBeNil)

					So(len(*query.Result), ShouldEqual, 2)
				})
			})
		})
	})
}