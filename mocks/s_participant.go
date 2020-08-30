package mocks

import (
	"github.com/sasalatart/batcoms/domain"
)

// SFaction returns a faction instance of domain.SParticipant that may be used for testing purposes
func SFaction() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.FactionKind,
		ID:          21418258,
		URL:         "https://en.wikipedia.org/wiki/French_First_Empire",
		Name:        "First French Empire",
		Description: "Empire of Napoleon I of France between 1804–1815",
		Extract:     "The First French Empire, officially the French Empire or the Napoleonic Empire, was the empire of Napoleon Bonaparte of France and the dominant power in much of continental Europe at the beginning of the 19th century. Although France had already established an overseas colonial empire beginning in the 17th century, the French state had remained a kingdom under the Bourbons and a republic after the French Revolution. Historians refer to Napoleon's regime as the First Empire to distinguish it from the restorationist Second Empire (1852–1870) ruled by his nephew Napoleon III.",
	}
}

// SFaction2 returns a faction instance of domain.SParticipant that may be used for testing purposes
func SFaction2() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.FactionKind,
		ID:          20611504,
		URL:         "https://en.wikipedia.org/wiki/Imperial_Russia",
		Name:        "Russian Empire",
		Description: "Empire in Eurasia and North America",
		Extract:     "The Russian Empire was an empire that extended across Eurasia and North America from 1721, following the end of the Great Northern War, until the Republic was proclaimed by the Provisional Government that took power after the February Revolution of 1917. The third-largest empire in history, at its greatest extent stretching over three continents, Europe, Asia, and North America, the Russian Empire was surpassed in size only by the British and Mongol empires. The rise of the Russian Empire coincided with the decline of neighboring rival powers: the Swedish Empire, the Polish–Lithuanian Commonwealth, Persia and the Ottoman Empire. It played a major role in 1812–1814 in defeating Napoleon's ambitions to control Europe and expanded to the west and south.",
	}
}

// SFaction3 returns a faction instance of domain.SParticipant that may be used for testing purposes
func SFaction3() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.FactionKind,
		ID:          266894,
		URL:         "https://en.wikipedia.org/wiki/Austrian_Empire",
		Name:        "Austrian Empire",
		Description: "monarchy in Central Europe between 1804 and 1867",
		Extract:     "The Austrian Empire was a Central European multinational great power from 1804 to 1867, created by proclamation out of the realms of the Habsburgs. During its existence, it was the third most populous empire after the Russian Empire and the United Kingdom in Europe. Along with Prussia, it was one of the two major powers of the German Confederation. Geographically, it was the third largest empire in Europe after the Russian Empire and the First French Empire. Proclaimed in response to the First French Empire, it partially overlapped with the Holy Roman Empire until the latter's dissolution in 1806.",
	}
}

// SCommander returns a commander instance of domain.SParticipant that may be used for testing purposes
func SCommander() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.CommanderKind,
		ID:          69880,
		URL:         "https://en.wikipedia.org/wiki/Emperor_Napoleon_I",
		Name:        "Napoleon",
		Description: "19th century French military leader, strategist, and politician",
		Extract:     "Napoleon Bonaparte, born Napoleone di Buonaparte, was a French statesman and military leader who became famous as an artillery commander during the French Revolution. He led many successful campaigns during the French Revolutionary Wars and was Emperor of the French as Napoleon I from 1804 until 1814 and again briefly in 1815 during the Hundred Days. Napoleon dominated European and global affairs for more than a decade while leading France against a series of coalitions during the Napoleonic Wars. He won many of these wars and a vast majority of his battles, building a large empire that ruled over much of continental Europe before its final collapse in 1815. He is considered one of the greatest commanders in history, and his wars and campaigns are studied at military schools worldwide. Napoleon's political and cultural legacy has made him one of the most celebrated and controversial leaders in human history.",
	}
}

// SCommander2 returns a commander instance of domain.SParticipant that may be used for testing purposes
func SCommander2() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.CommanderKind,
		ID:          27126603,
		URL:         "https://en.wikipedia.org/wiki/Alexander_I_of_Russia",
		Name:        "Alexander I of Russia",
		Description: "Emperor of Russia",
		Extract:     "Alexander I was the Emperor of Russia (Tsar) between 1801 and 1825. He was the eldest son of Paul I and Sophie Dorothea of Württemberg. Alexander was the first king of Congress Poland, reigning from 1815 to 1825, as well as the first Russian Grand Duke of Finland, reigning from 1809 to 1825.",
	}
}

// SCommander3 returns a commander instance of domain.SParticipant that may be used for testing purposes
func SCommander3() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.CommanderKind,
		ID:          251000,
		URL:         "https://en.wikipedia.org/wiki/Mikhail_Illarionovich_Kutuzov",
		Name:        "Mikhail Kutuzov",
		Description: "Field Marshal of the Russian Empire",
		Extract:     "Prince Mikhail Illarionovich Golenishchev-Kutuzov was a Field Marshal of the Russian Empire. He served as one of the finest military officers and diplomats of Russia under the reign of three Romanov Tsars: Catherine II, Paul I and Alexander I. His military career was closely associated with the rising period of Russia from the end of the 18th century to the beginning of the 19th century. Kutuzov is considered to have been one of the best Russian generals.",
	}
}

// SCommander4 returns a commander instance of domain.SParticipant that may be used for testing purposes
func SCommander4() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.CommanderKind,
		ID:          11551,
		URL:         "https://en.wikipedia.org/wiki/Francis_II,_Holy_Roman_Emperor",
		Name:        "Francis II, Holy Roman Emperor",
		Description: "The last Holy Roman Emperor and first Emperor of Austria",
		Extract:     "Francis II was the last Holy Roman Emperor, ruling from 1792 until 6 August 1806, when he dissolved the Holy Roman Empire after the decisive defeat at the hands of the First French Empire led by Napoleon at the Battle of Austerlitz. In 1804, he had founded the Austrian Empire and became Francis I, the first Emperor of Austria, ruling from 1804 to 1835, so later he was named the first Doppelkaiser in history.. For the two years between 1804 and 1806, Francis used the title and style by the Grace of God elected Roman Emperor, ever Augustus, hereditary Emperor of Austria and he was called the Emperor of both the Holy Roman Empire and Austria. He was also Apostolic King of Hungary, Croatia and Bohemia as Francis I. He also served as the first president of the German Confederation following its establishment in 1815.",
	}
}

// SCommander5 returns a commander instance of domain.SParticipant that may be used for testing purposes
func SCommander5() domain.SParticipant {
	return domain.SParticipant{
		Kind:        domain.CommanderKind,
		ID:          14092123,
		URL:         "https://en.wikipedia.org/wiki/Franz_von_Weyrother",
		Name:        "Franz von Weyrother",
		Description: "Austrian general",
		Extract:     "Franz von Weyrother was an Austrian staff officer and general who fought during the French Revolutionary Wars and the Napoleonic Wars. He drew up the plans for the disastrous defeats at the Battle of Rivoli, Battle of Hohenlinden and the Battle of Austerlitz, in which the Austrian army was defeated by Napoleon Bonaparte twice and Jean Moreau once.",
	}
}
