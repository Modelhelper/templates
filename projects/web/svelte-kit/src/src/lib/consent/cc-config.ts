import { COOKIE_CONSENT_NAME } from '$lib/constants';
import * as CookieConsent from 'vanilla-cookieconsent';

const config: CookieConsent.CookieConsentConfig = {
	cookie: {
		name: COOKIE_CONSENT_NAME,
	},
	categories: {
		necessary: {
			enabled: true,
			readOnly: true
		},
		analytics: {
			autoClear: {
				cookies: [
					{
						name: /^_ga/ // regex: match all cookies starting with '_ga'
					},
					{
						name: '_gid' // string: exact cookie name
					}
				]
			},

			// https://cookieconsent.orestbida.com/reference/configuration-reference.html#category-services
			services: {
				ga: {
					label: 'Google Analytics',
					onAccept: () => {},
					onReject: () => {}
				},
				// another: {
				// 	label: 'Another service',
				// 	onAccept: () => {},
				// 	onReject: () => {}
				// }
			}
		},
		ads: {}
	},

	onFirstConsent: ({ cookie }) => {
		console.log('onFirstConsent fired', cookie);
	},

	onConsent: ({ cookie }) => {
		console.log('onConsent fired!', cookie, CookieConsent.getUserPreferences());
	},

	onChange: ({ changedCategories, changedServices }) => {
		console.log('onChange fired!', changedCategories, changedServices);
	},

	onModalReady: ({ modalName }) => {
		console.log('ready:', modalName);
	},

	onModalShow: ({ modalName }) => {
		console.log('visible:', modalName);
	},

	onModalHide: ({ modalName }) => {
		console.log('hidden:', modalName);
	},

	guiOptions: {
		consentModal: {
			layout: 'cloud',
			position: 'bottom center',
			equalWeightButtons: true,
			flipButtons: false
		},
		preferencesModal: {
			layout: 'box',
			equalWeightButtons: true,
			flipButtons: false
		},
        
	},
    disablePageInteraction: false,
	language: {
		default: 'en',
        // autoDetect: "browser",
		translations: {
			en: {
				consentModal: {
					title: 'We use cookies',
					description:
						'We use cookies to store information so that the website is adapted to your choices and you get access to protected pages. Some functions on our websites need to store information in your browser in order to function.',
					acceptAllBtn: 'Accept',
					acceptNecessaryBtn: 'Reject all',
					showPreferencesBtn: 'Manage Individual preferences',
                    
					closeIconLabel: 'Reject all and close modal',
					// footer: `
					// 		<a href="#path-to-impressum.html" target="_blank">Impressum</a>
					// 		<a href="#path-to-privacy-policy.html" target="_blank">Privacy Policy</a>
					// `
				},
				preferencesModal: {
					title: 'Manage cookie preferences',
					acceptAllBtn: 'Accept all',
					acceptNecessaryBtn: 'Reject all',
					savePreferencesBtn: 'Accept current selection',
					closeIconLabel: 'Close modal',
					serviceCounterLabel: 'Service|Services',
					sections: [
						{
							title: 'Your Privacy Choices',
							description: `In this panel you can express some preferences related to the processing of your personal information. You may review and change expressed choices at any time by resurfacing this panel via the provided link. To deny your consent to the specific processing activities described below, switch the toggles to off or use the “Reject all” button and confirm you want to save your choices.`
						},
						{
							title: 'Strictly Necessary',
							description:
								'These cookies are essential for the proper functioning of the website and cannot be disabled. Without them you cannot log into the system or manage your privacy choices.',

							//this field will generate a toggle linked to the 'necessary' category
							linkedCategory: 'necessary'
						},
						{
							title: 'Performance and Analytics',
							description:
								'These cookies collect information about how you use our website. All of the data is anonymized and cannot be used to identify you.',
							linkedCategory: 'analytics',
							cookieTable: {
								caption: 'Cookie table',
								headers: {
									name: 'Cookie',
									domain: 'Domain',
									desc: 'Description'
								},
								body: [
									{
										name: '_ga',
										domain: 'yourdomain.com',
										desc: 'Description 1'
									},
									{
										name: '_gid',
										domain: 'yourdomain.com',
										desc: 'Description 2'
									}
								]
							}
						},
						{
							title: 'Targeting and Advertising',
							description:
								'These cookies are used to make advertising messages more relevant to you and your interests. The intention is to display ads that are relevant and engaging for the individual user and thereby more valuable for publishers and third party advertisers.',
							linkedCategory: 'ads'
						},
						{
							title: 'More information',
							description:
								'For any queries in relation to our policy on cookies and your choices, please <a href="#contact-page">contact us</a>'
						}
					]
				}
			},
			no: {
				consentModal: {
					title: 'Vi bruker informasjonskapsler',
					description:
						'Vi bruker informasjonskapsler for å lagre informasjon slik at nettstedet tilpasses dine valg og du får tilgang til beskyttede sider. Noen funksjoner på nettstedene våre trenger å lagre informasjon i nettleseren din for å fungere.',
					acceptAllBtn: 'Godta',
					acceptNecessaryBtn: 'Avvis alle',
					showPreferencesBtn: 'Administrer Individuelle preferanser',
                    
					closeIconLabel: 'Avvis alle og lukk dialog',
					// footer: `
					// 		<a href="#path-to-impressum.html" target="_blank">Impressum</a>
					// 		<a href="#path-to-privacy-policy.html" target="_blank">Privacy Policy</a>
					// `
				},
				preferencesModal: {
					title: 'Sett informasjonskapsler',
					acceptAllBtn: 'Godta alle',
					acceptNecessaryBtn: 'Avvis alle',
					savePreferencesBtn: 'Godta valgt',
					closeIconLabel: 'Lukk',
					serviceCounterLabel: 'Tjeneste|Tjenester',
					sections: [
						{
							title: 'Dine personvernvalg',
							description: `I dette panelet kan du sette noen preferanser knyttet til behandlingen av dine personopplysninger. Du kan gjennomgå og endre uttrykte valg når som helst ved å gjenoppta dette panelet via den angitte lenken. For å nekte samtykke til de spesifikke behandlingsaktivitetene som er beskrevet nedenfor, slår du av bryterne eller bruker knappen "Avvis alle" og bekrefter at du vil lagre valgene dine.`
						},
						{
							title: 'Nødvendig',
							description:
								'Disse informasjonskapslene er nødvendige for at nettstedet skal fungere riktig og kan ikke deaktiveres. Uten dem kan du ikke logge inn på systemet eller administrere personvernvalgene dine.',

							//this field will generate a toggle linked to the 'necessary' category
							linkedCategory: 'necessary'
						},
						{
							title: 'Statistiske',
							description:
								'Disse informasjonskapslene samler inn informasjon om hvordan du bruker nettstedet vårt. Alle dataene er anonymiserte og kan ikke brukes til å identifisere deg.',
							linkedCategory: 'analytics',
							cookieTable: {
								caption: 'Cookie- tabell',
								headers: {
									name: 'Navn',
									domain: 'Domene',
									desc: 'Beskrivelse'
								},
								body: [
									{
										name: '_ga',
										domain: 'ocean-training.no',
										desc: 'Hovedinformasjonskapsel for Google Analytics'
									},
									{
										name: '_gid',
										domain: 'ocean-training.no',
										desc: 'Informasjonskapsel for Google Analytics som lagrer og oppdaterer en unik verdi for hver side du besøker'
									}
								]
							}
						},						
						{
							title: 'Mer informasjon',
							description:
								'For spørsmål knyttet til vår retningslinje for informasjonskapsler og dine valg, vennligst <a href="#contact-page">kontakt oss</a>'
						}
					]
				}
			}
		}
	}
};

export default config;
