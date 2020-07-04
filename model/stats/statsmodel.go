/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 *
 */

package statsModel

import (
    "fmt"
    "strconv"

    "github.com/jmoiron/sqlx"
)


type Model struct {
    dbx *sqlx.DB
}

type TableMeasure struct {
    CompanyName             string   `db:"company_name"             json:"companyName"`
    MetrixPeriod            string   `db:"metrix_period"            json:"metrixPeriod"`
    IntegralValue           float32  `db:"integral_value"           json:"integralValue"`
    FirstprocessingValue    float32  `db:"firstprocessing_value"    json:"firstprocessingValue"`
    SecondprocessingValue   float32  `db:"secondprocessing_value"   json:"secondprocessingValue"`
    LifecycleValue          float32  `db:"lifecycle_value"          json:"lifecycleValue"`
    AppealcountValue        float32  `db:"appealcount_value"        json:"appealcountValue"`
}

type SLMeasure struct {
    Label   string          `json:"label"`
    Value   float32         `json:"value"`
}

type Company struct {
    FullName    string      `json:"fullName"`
    SLData      []SLMeasure `json:"slData"`
}

type StatsBatch struct {
    ComsByIntegral          []Company   `json:"comsByIntegral"`
    ComsByFirstProcessing   []Company   `json:"comsByFirstProcessing"`
    ComsBySecondProcessing  []Company   `json:"comsBySecondProcessing"`
    ComsByLifecicle         []Company   `json:"comsByLifecicle"`
    ComsByAppealcount       []Company   `json:"comsByAppealcount"`
}

func (this *Model) GetStats(mspName string, year string) (StatsBatch, error) {

    var request string
    var err error
    var statsBatch StatsBatch

    statsBatch.ComsByIntegral =           make([]Company, 0)
    statsBatch.ComsByFirstProcessing =    make([]Company, 0)
    statsBatch.ComsBySecondProcessing =   make([]Company, 0)
    statsBatch.ComsByLifecicle =          make([]Company, 0)
    statsBatch.ComsByAppealcount =        make([]Company, 0)

    yearInt, err := strconv.Atoi(year)
    if err != nil {
        return statsBatch, err
    }

    request = `
        SELECT DISTINCT company_name
        FROM arrangement.get_claims_sla_metrix(
                '%d-01-01'::DATE,
                '%d-01-01'::DATE,
                $1
        )
        ORDER BY company_name`

    request = fmt.Sprintf(request, yearInt, yearInt + 1)

    var companyNames []string
    err = this.dbx.Select(&companyNames, request, mspName)
    if err != nil {
        return statsBatch, err
    }

    request = `
        SELECT
            company_name,
            metrix_period,
            integral_value,
            firstprocessing_value,
            secondprocessing_value,
            lifecycle_value,
            appealcount_value
        FROM arrangement.get_claims_sla_metrix(
                '%d-01-01'::DATE,
                '%d-01-01'::DATE,
                $1
        )
        ORDER BY company_name, metrix_period`

    request = fmt.Sprintf(request, yearInt, yearInt + 1)

    var measures []TableMeasure
    err = this.dbx.Select(&measures, request, mspName)
    if err != nil {
        return statsBatch, err
    }

    for _, companyName := range companyNames {
        tmpComByIntegral := Company{
            FullName: companyName,
        }
        tmpComByFirstProcessing := Company{
            FullName: companyName,
        }
        tmpComBySecondProcessing := Company{
            FullName: companyName,
        }
        tmpComByLifecicle := Company{
            FullName: companyName,
        }
        tmpComByAppealcount := Company{
            FullName: companyName,
        }

        for _, measure := range measures {
            if measure.CompanyName == companyName {
                tmpComByIntegral.SLData = append(tmpComByIntegral.SLData, SLMeasure{
                        Label: measure.MetrixPeriod,
                        Value: measure.IntegralValue,
                    })
                tmpComByFirstProcessing.SLData = append(tmpComByFirstProcessing.SLData, SLMeasure{
                        Label: measure.MetrixPeriod,
                        Value: measure.FirstprocessingValue,
                    })
                tmpComBySecondProcessing.SLData = append(tmpComBySecondProcessing.SLData, SLMeasure{
                        Label: measure.MetrixPeriod,
                        Value: measure.SecondprocessingValue,
                    })
                tmpComByLifecicle.SLData = append(tmpComByLifecicle.SLData, SLMeasure{
                        Label: measure.MetrixPeriod,
                        Value: measure.LifecycleValue,
                    })
                tmpComByAppealcount.SLData = append(tmpComByAppealcount.SLData, SLMeasure{
                        Label: measure.MetrixPeriod,
                        Value: measure.AppealcountValue,
                    })
            }
        }
        statsBatch.ComsByIntegral =           append(statsBatch.ComsByIntegral, tmpComByIntegral)
        statsBatch.ComsByFirstProcessing =    append(statsBatch.ComsByFirstProcessing, tmpComByFirstProcessing)
        statsBatch.ComsBySecondProcessing =   append(statsBatch.ComsBySecondProcessing, tmpComBySecondProcessing)
        statsBatch.ComsByLifecicle =          append(statsBatch.ComsByLifecicle, tmpComByLifecicle)
        statsBatch.ComsByAppealcount =        append(statsBatch.ComsByAppealcount, tmpComByAppealcount)
    }

    return statsBatch, nil
}

func New(dbx *sqlx.DB) *Model {
    return &Model{
        dbx: dbx,
    }
}
