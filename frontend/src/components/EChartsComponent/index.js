import * as echarts from 'echarts';
import { useEffect, useRef } from "react"


export default function EChartsComponent(props) {
    
    const chartsRef = useRef();
    const { 
        title, 
        xAxis, 
        yAxis, 
        series, 
        tooltip,
        width,
        height
    } = props;

    function init() {
        const chart = echarts.init(chartsRef.current);
        chart.setOption({
            title: {
                text: title
            },
            tooltip,
            xAxis,
            yAxis,
            series
        })
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div 
            className="chartsContainer" 
            ref={chartsRef}
            style={{
                width: width,
                height: height
            }}
        >

        </div>
    )
}