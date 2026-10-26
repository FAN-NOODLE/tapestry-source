import 'axios';

declare module 'axios' {
  interface AxiosRequestConfig {
    upload?: boolean;
    download?: boolean;
    hideLoading?: boolean;
  }
}
