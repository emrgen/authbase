import dayjs from "dayjs";
import {useCallback} from "react";
import {useNavigate} from "react-router";
import {authbase} from "../api/client.ts";

let interval: string | number | NodeJS.Timeout | null | undefined = null;

export function useRotateAccessToken() {
  const navigate = useNavigate();

  return useCallback(() => {
    clearInterval(interval as NodeJS.Timeout);

    interval = setInterval(() => {
      const expiresAt = localStorage.getItem("expiresAt");
      const expiresAtDate = dayjs(expiresAt).toDate();
      const now = new Date();
      const diff = expiresAtDate.getTime() - now.getTime();
      // if token expires in less than 5 minutes, refresh it
      if (diff > 5 * 60 * 1000) {
        return;
      }

      const refreshToken = localStorage.getItem("refreshToken");
      if (!refreshToken) {
        return Promise.reject(new Error("No refresh token found"));
      }
      return authbase.auth.refresh({
        body: {
          refreshToken,
        },
      }).then((res) => {
        const {data} = res;
        const {tokens = {},} = data;
        const {accessToken = '', refreshToken = '', expiresAt = ''} = tokens;
        localStorage.setItem("accessToken", accessToken.toString());
        localStorage.setItem("refreshToken", refreshToken.toString());
        localStorage.setItem("expiresAt", expiresAt.toString());
      }).catch((err) => {
        console.log(err);
        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
        localStorage.removeItem("expiresAt");
        clearInterval(interval as NodeJS.Timeout);
        navigate('/login');
      })
    }, 1000); // refresh the token 5 minutes before it expires
  }, [navigate]);
}