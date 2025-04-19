import {authbase} from "../api/client.ts";

let interval = null;

export function rotateAccessToken() {
  clearInterval(interval);
  interval = setInterval(() => {
    const refreshToken = localStorage.getItem("refreshToken");
    if (!refreshToken) {
      return Promise.reject(new Error("No refresh token found"));
    }
    return authbase.auth.refresh({
      body: {
        refreshToken,
      },
    }).then((res) => {
      console.log(res);
      // localStorage.setItem("accessToken", accessToken.toString());
      // localStorage.setItem("refreshToken", refreshToken.toString());
    });
  }, 1000 * 60 * 5); // 5 minutes
}