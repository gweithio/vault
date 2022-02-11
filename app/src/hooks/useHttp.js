const defaultHeaders = {
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin':'*',
};

const defaultUrl = 'http://127.0.0.1:8080';

export const useGet = async (url, params = {}) => {
    return fetch(defaultUrl + url, {
        method: "GET",
        mode: 'no-cors',
        headers: defaultHeaders,
    }).then(res => res.json()).then(res => res);
}

export const usePost = async (url, params = {}) => {
    return fetch(defaultUrl + url, {
        method: "POST",
        mode: 'no-cors',
        headers: defaultHeaders,
        body: JSON.stringify(params)
    }).then(res => res.json()).then(res => res);

}

export const useDelete = async (url, params = {}) => {
    return fetch(defaultUrl + url, {
        method: "DELETE",
        mode: 'no-cors',
        headers: defaultHeaders,
        body: JSON.stringify(params)
    }).then(res => res.json()).then(res => res);
}