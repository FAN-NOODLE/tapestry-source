export class AbortAblePromise {
  promise: Promise<any>;
  aborted: boolean;
  constructor(promise: Promise<any>) {
    this.promise = promise;
    this.aborted = false;
  }

  abort() {
    this.aborted = true;
  }

  then(onFulfilled?: (value: any) => any, onRejected?: (value: any) => any) {
    if (this.aborted) {
      return Promise.reject(new Error('Promise aborted'));
    }
    return this.promise.then(onFulfilled, onRejected);
  }

  catch(onRejected: (value: any) => any) {
    if (this.aborted) {
      return Promise.reject(new Error('Promise aborted'));
    }
    return this.promise.catch(onRejected);
  }

  finally(onFinally: () => void) {
    if (this.aborted) {
      return Promise.reject(new Error('Promise aborted'));
    }
    return this.promise.finally(onFinally);
  }
}