import { userApi } from "../user";
import Request from "./index";
import axios, { HttpStatusCode, type AxiosError } from "axios";
import type { RequestInterceptors } from "./type";
import { useLocalStorage } from "@vueuse/core";
const AccessToken = useLocalStorage("access_token", '');
const RefreshToken = useLocalStorage("refresh_token", '');
export const _RequstInterceptors: RequestInterceptors = {
  requestInterceptors(config) {
    config.headers["Authorization"] ||= `Bearer ${AccessToken.value}`;
    return config;
  },
  requestInterceptorsCatch(err) {
    return err;
  },
  responseInterceptors(res) {
    if (res.data.code !== HttpStatusCode.Unauthorized) {
      return Promise.resolve(res);
    }

    // 清除令牌并跳转到登录页
    const clearAndRedirect = (msg?: string) => {
      console.error("清除令牌并提示登录", msg);
      const btn = window.confirm(msg || "登录已过期，请重新登录");
      if (btn) {
        window.location.href = "/user/auth";
      }
      // TODO调用消息
      [AccessToken.value, RefreshToken.value] = [null, null];
      return Promise.reject(res);
    };

    // 没有刷新令牌，直接清除并跳转
    if (!RefreshToken.value || RefreshToken.value == "") {
      return clearAndRedirect("没有令牌, 请重新登录");
    }

    // 尝试刷新令牌
    console.error("登录已过期，正在尝试刷新令牌... ");
    return userApi
      .RefreshToken()
      .then(({ data, config, status }) => {
        const { code, info, msg } = data;
        if (
          code !== HttpStatusCode.Ok ||
          info == undefined ||
          info.length !== 2
        ) {
          console.error("刷新令牌失败-1", msg);
          return clearAndRedirect(msg);
        }

        // 刷新成功，重试请求
        [AccessToken.value, RefreshToken.value] = [info[0], info[1]];
        return Request.request({
          ...res.config,
          headers: {
            ...res.config.headers,
            Authorization: `Bearer ${AccessToken.value}`,
          },
        });
      })
      .catch(() => {
        console.error("刷新令牌失败-2");
        return clearAndRedirect("刷新令牌失败");
      });
  },
  responseInterceptorsCatch(axiosInstance, err: AxiosError) {
    const message = err.code === "服务端主动终止连接" ? "请求超时" : undefined;
    if (axios.isCancel(err)) {
      return Promise.reject(err);
    }
    // TODO: 处理返回状态
    // checkErrorStatus(
    //     (err as AxiosError).response?.status,
    //     message,
    //     (message) =>
    //         console.log(message)
    // )
    // return retry(axiosInstance, err as AxiosError)
  },
};
