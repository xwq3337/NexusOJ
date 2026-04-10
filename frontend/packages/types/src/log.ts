export interface Log {
    level : string // INFO | WARN | ERROR
    ts : string // timestamp 2026-03-12T16:39:03.559+0800
    caller : string // caller function name
    msg : string // log message
    status : number // http status code
    method : string // http method
    path : string // http path
    ip : string // client ip
    user_id? : string // user id
    latency : string // request latency 335.791µs
    user_agent : string // client user agent
    headers : Object  // 
    response_body : Object //
}