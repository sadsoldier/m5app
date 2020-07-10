# claimstat-api


### pre-Form requests

- Тип страхования, все возможные значения 

```
    SELECT DISTINCT insurance_type FROM flex.claims;
     insurance_type 
    ----------------
     Неизвестно     
     Грузы          
     Сроки
     (3 rows)
```


- Статус претензии, все возможные значения 

```
    SELECT DISTINCT status FROM flex.claims;
     status                      
    -----------------------------
     Новая претензия             
     Принята                     
     Рассмотрение страховщиком   
     Корректировка страхователем 
     Отозвана                    
     Отказано                    
     Удовлетворена               
    (7 rows)
```


- Список страхователей (insurers) данной ТК 

```
    SELECT DISTINCT insurance_company_msp, insurance_company FROM flex.claims WHERE transport_company_msp = 'DellinMSP';
     insurance_company_msp | insurance_company                 
    -----------------------+-----------------------------------
     AlfastrahMSP          | АО «АльфаСтрахование»             
     IngosMSP              | СПАО «Ингосстрах»                 
     RenaissanceMSP        | AО «Группа Ренессанс Страхование» 
    (3 rows)
```


- Список страховщиков (policyHolders) данной СК

```
    SELECT DISTINCT transport_company_msp, transport_company FROM flex.claims WHERE insurance_company_msp = 'AlfastrahMSP';
     transport_company_msp | transport_company         
    -----------------------+---------------------------
     MisMSP                | ООО «МИС»                 
     SkifMSP               | ООО «Компания Скиф-Карго» 
     DellinMSP             | ООО «Деловые линии»       
     AvtoTransitMSP        | ООО «ЦАП 2015»            
    (4 rows)
```


- Валюта, все возможные значения
```
    SELECT DISTINCT currency FROM flex.claims;
     currency 
    ----------
     rub      
    (1 row)
```


