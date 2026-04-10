export type ApiResponse<T = any> =  {
  code: number;
  msg: string;
  info?: T;
}

export type VerdictType = "Accepted" | "WrongAnswer" | "TimeLimitExceeded" | "MemoryLimitExceeded" | "RuntimeError" | "CompilationError"


export * from './src/record';

export * from './src/chat';

export * from './src/user';

export * from './src/blog';

export * from './src/contest';

export * from './src/problem';

export * from './src/log'

export * from './src/solution'