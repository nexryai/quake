export const getQuakeScaleColor = (scale: number): string => {
    switch (scale) {
        case 10: return "#b3d5ff"; // 震度1
        case 20: return "#00c980"; // 震度2
        case 30: return "#ffb630"; // 震度3
        case 40: return "#ff7926"; // 震度4
        case 45: return "#fc0505"; // 震度5弱
        case 50: return "#fc0505"; // 震度5強
        case 55: return "#820101"; // 震度6弱 (赤)
        case 60: return "#820101"; // 震度6強 (濃い赤)
        case 70: return "#c400d6"; // 震度7 (紫)
        default: return "#CCCCCC"; // 不明な値 (グレー)
    }
};

export const getQuakeScaleColorText = (scale: number): string => {
    return scale <= 40 ? "black" : "white";
};
