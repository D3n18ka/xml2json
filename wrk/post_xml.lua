local f = io.open("data.xml", "rb")
wrk.body   = f:read("*all")
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/xml"
