import {
  AdminAuthServiceApi,
  AuthServiceApi,
  Configuration,
  ProjectServiceApi,
  AccountServiceApi,
  PoolServiceApi,
} from '@emrgen/authbase-client-gen';
import {AxiosInstance} from "axios";

export class Config {
  token: () => string;
  basePath: string;
  axios?: AxiosInstance;

  constructor(token: () => string, basePath: string, axios?: AxiosInstance) {
    this.token = token;
    this.basePath = basePath;
    this.axios = axios;
  }
}

// AuthbaseClient is a wrapper around the Authbase API client
export class AuthbaseClient {
  auth: AuthServiceApi;
  account: AccountServiceApi;
  project: ProjectServiceApi;
  pool: PoolServiceApi;
  // member: MemberServiceApi;
  // token: OfflineTokenServiceApi;
  admin: AdminAuthServiceApi

  constructor(config: Configuration) {


    this.account = new AccountServiceApi(config);
    this.project = new ProjectServiceApi(config);
    this.pool = new PoolServiceApi(config);
    // this.member = new MemberServiceApi(config);
    // this.token = new OfflineTokenServiceApi(config);
    this.auth = new AuthServiceApi(config);
    this.admin = new AdminAuthServiceApi(config);
  }
}


export class TokenConfig extends Configuration {
  constructor() {
    super({
      accessToken: () => {
        return localStorage.getItem('accessToken') ?? '---';
      },
      basePath: 'http://localhost:4001', // process.env.REACT_APP_AUTHBASE_API_URL || 'http://localhost:8080',
    });
  }
}


const config = new TokenConfig();

export const authbase = new AuthbaseClient(config);