import axios from 'axios'
import { createClient } from 'redis'

type Payload = {
    industry: string;
    Comp_domains: string[] | null;
    browser: string;
}

export async function getForumProductLinks(payload?: Payload) {
    try {
        // const res = await axios.post("http://localhost:8080/back-link",{"comp_domains":["amazon.com","ajmadison.com","homedepot.com","bestbuy.com","build.com","lowes.com"],"industry":`${payload?.industry}`,"browser":"duckduckgo"},{'headers':{'Content-Type' : 'application/json'}})
        const res = await axios.post("http://localhost:8080/forum-scrape",{"industry":`${payload?.industry}`,"browser":"duckduckgo","comp_domains":payload?.Comp_domains},{'headers':{'Content-Type' : 'application/json'}, withCredentials: true})
        const data = await res.data
        return data;
    } catch (error: any) {
        console.log("error",error)
    }
}

export async function getBacklink(id: string) {
    const res = await axios.post(`http://localhost:8080/back-link/${id}`);
    const data = await res.data;
    return data;
}

export async function getBacklinks() {
    const res = await axios.get('http://localhost:8080/api/user/backlinks', {withCredentials: true})
    const data = await res.data;
    return data
}

export async function getGenAIReply(url: string) {
    const queryParam = encodeURI(url)
    const res = await axios.get(`http://localhost:8000/reply?url=${queryParam}`)
    const data = await res.data;
    return data
}

export async function testGenAiReply() : Promise<any> {
    const res = "Great thread! I completely understand your journey - I also went through that same realization a couple of years ago when I found myself constantly adjusting my larger watches under winter layers. That 11.4mm thickness on the Helson Shark Diver that badgerracer mentioned is absolutely insane for a 300m watch!\n\nYou've got some fantastic suggestions here already. The Lorier Neptune keeps coming up for good reason - it really does punch above its weight class. I noticed several members mentioned the Christopher Ward Trident Pro 38mm, and I have to agree it's a solid choice when you catch it on sale.\n\nOne option that might be perfect for your needs is the **[PRODUCT LINK: https://www.example-watches.com/seiko-spb143-38mm-diver]** which hits that sweet spot of being under 39mm while maintaining excellent proportions. What I love about it is how it manages to feel substantial on the wrist without being bulky under cuffs. The lug-to-lug is compact enough that it should work beautifully with your wrist size.\n\nThe lug-to-lug measurement that mconlonx mentioned is absolutely crucial - sometimes a watch with a smaller diameter but longer lug-to-lug can wear larger than its specs suggest. Have you considered trying any of the vintage options that Cvp33 suggested? Some of those older divers have incredible character and proportions that modern watches often miss.\n\nKeep us posted on what you decide - I'm always excited to see fellow enthusiasts find that perfect fitting piece!"
    return new Promise(resolve => {
        setTimeout(() => {
            resolve(res)
        },1000)
    })
}

export async function publishBacklink(id: number, response: string) {
    const res = await axios.post(`http://localhost:8080/api/user/backlink/${id}`, {response: response}, {withCredentials:true})
    const data = await res.data
    return data;
}