:: --------------------------------------------------------------
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

echo.|edgepulse.exe -o "result_ddns.txt"

if not exist result_ddns.txt (
    echo.
    echo EdgePulse: operation completed.
    goto :END
)

for /f "skip=1 tokens=1 delims=," %%i in (result_ddns.txt) do (
    Echo %%i
    if "%%i"=="" (
        echo.
        echo EdgePulse: operation completed.
        goto :END
    )
    curl -X PUT "https://api.cloudflare.com/client/v4/zones/ID/dns_records/ID" ^
            -H "X-Auth-Email: " ^
            -H "X-Auth-Key:  API " ^
            -H "Content-Type: application/json" ^
            --data "{\"type\":\"A\",\"name\":\"\",\"content\":\"%%i\",\"ttl\":1,\"proxied\":true}"
::            -H "Content-Type: application/json" ^

        goto :END
)
:END
pause