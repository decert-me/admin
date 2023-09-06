import { useRequest } from "ahooks";
import { Spin } from "antd";
import { useEffect } from "react";



export default function InfiniteScroll(props) {
    
    const { 
        className,              //  loading类名
        func,                   //  调用的方法
        debounceWait,           //  防抖的时间
        height,                 //  window 触发方法的高度
        component,              //  组件

        scrollRef,              //  scroll盒子dom
        isCustom,               //  结合scrollRef使用，true则是给scrollRef注册scroll事件, false则是给window注册
        customHeight,           //  scrollRef 触发方法的高度
    } = props;

    const { run } = useRequest(func, {
        debounceWait: 500,
        manual: true
    });

    function handleScroll(params) {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const windowHeight = window.innerHeight;
        const documentHeight = document.querySelector("main").scrollHeight;
        const zoom = document.body.style.zoom ? document.body.style.zoom : 1;
        if ((scrollTop + windowHeight + (height ? height : 50)) / zoom >= documentHeight) {
            run();
        }
    }

    function customHandleScroll(params) {
        const { scrollTop, clientHeight, scrollHeight } = scrollRef.current;
        const isLoading = document.querySelector(`${className ? '.'+className : ".loading"}`);
        if ((scrollTop + clientHeight >= (scrollHeight - (customHeight ? customHeight : 0))) && isLoading) {
            run();
        }
    }

    useEffect(() => {
        if (isCustom) {
            // ref监听
            scrollRef.current?.addEventListener('scroll', customHandleScroll);
            return () => {
                scrollRef.current?.removeEventListener('scroll', customHandleScroll);
            };
        }else{
            // window监听
            window.addEventListener('scroll', handleScroll);
            return () => {
                window.removeEventListener('scroll', handleScroll);
            };
        }
    }, []);


    return (
        component ?
        component
        :
        <Spin 
            size="large" 
            className={
                `${className ? className : "loading"}`
            } 
        />
    )
}