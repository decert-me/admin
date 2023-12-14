import Icon from '@ant-design/icons';


export default function CustomIcon(props) {

    const { type, className } = props;

    const star = (
      <svg t="1701682936461" className={className||`icon`} viewBox="0 0 1024 1024" fill="currentColor" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="9699" width="1em" height="1em">
        <path d="M956 398.496q-8-23.488-26.496-39.008t-42.496-19.488l-204.992-31.008-92-195.008q-11.008-24-32.992-36.992Q536.032 64 512.032 64t-44.992 12.992q-22.016 12.992-32.992 36.992l-92 195.008-204.992 31.008q-24 4-42.496 19.488t-26.496 39.008-2.496 47.008 22.496 41.504l151.008 154.016-36 218.016q-6.016 40 20 70.496t66.016 30.496q22.016 0 42.016-11.008l180.992-100 180.992 100q20 11.008 42.016 11.008 40 0 66.016-30.496t20-70.496l-36-218.016 151.008-154.016q16.992-18.016 22.496-41.504t-2.496-47.008z" p-id="4435"></path>
      </svg>
    )
    

    const icons = {
      "icon-star": star
    }

    return(
      <Icon component={() => icons[type]} />
    )
}