export interface Record {
    id : number
    problem_id : number
    user_id : string
    max_time : number
    max_memory : number
    code : string
    language : string
    created_at : string
    verdict : JudgeVerdictType
    updated_at : string
    judge_result : JudgeTestCaseResult[]
}

export interface JudgeTestCaseResult {
    case_id : string
    stdin : string
    stdout : string
    stderr : string
    status : JudgeVerdictType
    expected : string
    time : number
    memory : number
}
export type JudgeVerdictType = keyof typeof JudgeVerdict

export enum JudgeVerdict {
    Accepted = "Accepted",
    WrongAnswer = "WrongAnswer",
    TimeLimitExceeded = "TimeLimitExceeded",
    MemoryLimitExceeded = "MemoryLimitExceeded",
    RuntimeError = "RuntimeError",
    CompilationError = "CompilationError",
    RestrictedSystemCall = "RestrictedSystemCall",
    SystemError = "SystemError",
}




export interface GetRecordListParams {
    page : number
    page_size : number
    search : string
    verdict : JudgeVerdictType
    language : string
}

export interface GetRecordListResponse extends Record {
    problem_title : string
    username : string
}
export interface GetRecordDetailResponse extends Record {
    problem_title : string
    username : string
}