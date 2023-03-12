package models
import (
	"database/sql"
	"errors"
	"time"
)

type Event struct {
	EventID int
	EventName string
	Content string
	Venue string
	Technicality bool
	EventDate time.Time
}

type EventModel struct {
	DB *sql.DB
} 


func (m *EventModel) Insert(eventname, content, venue string, technicality bool, eventdate time.Time) (int, error) {
	var id int
	err:= m.DB.QueryRow(`INSERT INTO events (eventname, content, venue, technicality, eventdate) VALUES ($1, $2, $3, $4,$5) RETURNING eventid`, eventname, content, venue, technicality, eventdate).Scan(&id)
	if err!=nil{
		return  0,err
	}

	return id, nil
}

func (m *EventModel) Update(eventid int, eventname, content, venue string, technicality bool, eventdate time.Time) (int, error) {
	err:= m.DB.QueryRow(`UPDATE events SET eventname=$2, content=$3, venue=$4, technicality=$5, eventdate=$6 WHERE eventid=$1 RETURNING eventid`,eventid, eventname, content, venue, technicality, eventdate).Scan(&eventid)
	if err!=nil{
		return  0,err
	}
	return eventid, nil
}

func (m *EventModel) Get(eventid int) (*Event, error) {
	s := &Event{}

	err:= m.DB.QueryRow("SELECT eventid, eventname, content, venue, technicality, eventdate FROM events WHERE eventid = $1", eventid).Scan(&s.EventID, &s.EventName, &s.Content, &s.Venue, &s.Technicality, &s.EventDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *EventModel) ViewAll() ([]*Event, error) {
	events:= []*Event{}
	rows, err:= m.DB.Query("SELECT eventid, eventname, content, venue, technicality, eventdate FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		s:=&Event{}
		err=rows.Scan(&s.EventID, &s.EventName, &s.Content, &s.Venue, &s.Technicality, &s.EventDate)
	if err != nil {
		return nil, err
	}
	events=append(events,s)
	}
	if err=rows.Err();err!=nil{
		return nil, err
	}

	return events, nil
}

