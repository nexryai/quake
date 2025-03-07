<template>
    <div class="about">
        <div class="content">
            <h2>About this app</h2>
            <p>気象庁から地震・津波に関する情報を取得し、表示するウェブアプリです。</p>
            <div class="copyright">
                <p>©2024-2025 nexryai </p>
                <p>災害情報/マップデータ: ©<a href='https://www.jma.go.jp/jma/index.html' target='_blank'>気象庁</a></p>
                <p>
                    マップデータ（遠地地震用）: ©<a href='https://openstreetmap.org' target='_blank'>OpenStreetMap</a>
                    Licensed under <a href="opendatacommons.org">Open Database License</a>
                </p>
            </div>

            <h3>利用上の注意・免責事項</h3>
            <p>本サービスに関する可用性、正確性、信頼性、完全性、適時性、およびその他のいかなる保証も、明示的または黙示的であるかを問わず、一切提供されません。</p>
            <p>切迫した危険がある時に、このアプリからの情報のみに依存しないでください。必ず他の信頼できるソース、公的機関やNHKからの情報と併せて参照するようにしてください。</p>
            <br>
            <p> 本サービスを利用したことによって発生したいかなる損害も本サービスの運用に関連する法人、集団あるいは個人、及び本ソフトウェアのライセンサーは一切の責任を負わないものとします。
                本サービスの運営者は、一切の予告なしに本サービスの提供を停止、変更、または終了する権利、及び任意のユーザーやそのIPアドレスからのアクセスをブロックする権利を留保します。
                本サービスを利用することで、本サービスに関しての訴訟を起こす権利（集団訴訟、及びそれに参加する権利を含む）、調停や仲裁を含む裁判官や陪審員に審理してもらう権利を永久に放棄したことになります。
            </p>
            <br>
            <p>
                This service complies with laws and regulations in Japan and does not adhere to the GDPR (General Data
                Protection Regulation) of the European Union or the CCPA (California Consumer Privacy Act) of the United
                States. By using this service, you agree to waive any rights based on these laws. <br>
                本サービスは日本国内の法令・規制に準拠しており、 EUのGDPR（一般データ保護規則）
                や米国カリフォルニア州のCCPAに準拠していません。本サービスを利用することで、これらの法律に基づく権利を放棄したことになります。

            </p>
        </div>
    </div>
</template>

<style>
.about {
    margin: auto;

    & .content {
        margin: 16px;

        & .copyright {
            margin: 32px;

            & p {
                margin: 12px 0px 12px 0;
            }
        }

        & h2 {
            margin: 32px 0 16px 0;
        }

        & h3 {
            margin: 30px 0 14px 0;
        }
    }
}
</style>
