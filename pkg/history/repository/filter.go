package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/bagus2x/tjiwi/pkg/history"
	"github.com/sirupsen/logrus"
)

func descendingFilter(params *history.Params) (string, []interface{}) {
	columns := strings.Builder{}
	values := make([]interface{}, 0)
	stringIndex := 1

	if params.Direction == "prev" {
		columns.WriteString(`
			WITH prev_mode AS (
				SELECT
					h.id as history_id, bp.id, bp.gsm, bp.width, bp.io, bp.material_number, bp.quantity, bp.location, h.storage_id, 
					p.id, p.photo, p.username, h.status, h.affected, h.created_at
				FROM
					History h
				JOIN
					Base_Paper bp
				ON
					h.base_paper_id = bp.id
				JOIN
					Profile p
				ON
					h.member_id = p.id
				WHERE
		`,
		)
		if params.Cursor != 0 {
			fmt.Fprintf(&columns, " h.id > %d ", params.Cursor)

		} else {
			fmt.Fprintf(&columns, " h.id > %d ", 0)
		}
	} else {
		columns.WriteString(`
			SELECT
				h.id, bp.id, bp.gsm, bp.width, bp.io, bp.material_number, bp.quantity, bp.location, h.storage_id, 
				p.id, p.photo, p.username, h.status, h.affected, h.created_at
			FROM
				History h
			JOIN
				Base_Paper bp
			ON
				h.base_paper_id = bp.id
			JOIN
				Profile p
			ON
				h.member_id = p.id
			WHERE
		`,
		)

		logrus.Error("cursor", params.Cursor)

		if params.Cursor != 0 {
			fmt.Fprintf(&columns, " h.id < %d ", params.Cursor)

		} else {
			fmt.Fprintf(&columns, " h.id < %d ", math.MaxInt32)
		}
	}

	if params.StorageID != 0 {
		fmt.Fprintf(&columns, " AND h.storage_id = %d ", params.StorageID)
	}

	if params.Status != "" {
		fmt.Fprintf(&columns, " AND h.status = $%d ", stringIndex)
		values = append(values, params.Status)
		stringIndex++
	}

	if params.StartDate != 0 {
		fmt.Fprintf(&columns, " AND h.created_at >= %d ", params.StartDate)
	}

	if params.EndDate != 0 {
		fmt.Fprintf(&columns, " AND h.created_at <= %d ", params.EndDate)
	}

	if params.Direction == "prev" {
		columns.WriteString(" ORDER BY h.id ASC ")
	} else {
		columns.WriteString(" ORDER BY h.id DESC ")
	}

	if params.Limit != 0 {
		fmt.Fprintf(&columns, " LIMIT %d ", params.Limit)
	} else {
		fmt.Fprintf(&columns, " LIMIT %d ", 25)
	}

	if params.Direction == "prev" {
		columns.WriteString(` ) SELECT * FROM prev_mode ORDER BY history_id DESC`)
	}

	logrus.Info(columns.String())

	return columns.String(), values
}
