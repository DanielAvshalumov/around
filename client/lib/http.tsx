import axios from 'axios';

class Http {

    getIntroBacklinks = async (url: string) => {
        return await axios.get(`/api/intro/backlinks?url=${url}`);
    }
}