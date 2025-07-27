import dayjs from 'dayjs';

export interface SMARTCriteria {
  specific: boolean;
  measurable: boolean;
  achievable: boolean;
  relevant: boolean;
  timeBound: boolean;
}

export interface SMARTValidationResult {
  isValid: boolean;
  score: number; // 0-100
  criteria: SMARTCriteria;
  suggestions: string[];
  warnings: string[];
}

export interface GoalData {
  title: string;
  description: string;
  category: string;
  deadline?: string;
  priority: string;
}

/**
 * Validate a goal against SMART criteria
 */
export function validateSMARTGoal(goalData: GoalData): SMARTValidationResult {
  const criteria: SMARTCriteria = {
    specific: false,
    measurable: false,
    achievable: false,
    relevant: false,
    timeBound: false
  };

  const suggestions: string[] = [];
  const warnings: string[] = [];

  // S - Specific
  criteria.specific = validateSpecific(goalData.title, goalData.description);
  if (!criteria.specific) {
    suggestions.push("Make your goal more specific by clearly defining what exactly you want to achieve");
  }

  // M - Measurable
  criteria.measurable = validateMeasurable(goalData.title, goalData.description);
  if (!criteria.measurable) {
    suggestions.push("Add measurable metrics like numbers, percentages, or specific outcomes");
  }

  // A - Achievable
  criteria.achievable = validateAchievable(goalData);
  if (!criteria.achievable) {
    warnings.push("Consider if this goal is realistic given your current situation and resources");
  }

  // R - Relevant
  criteria.relevant = validateRelevant(goalData);
  if (!criteria.relevant) {
    suggestions.push("Ensure this goal aligns with your values and long-term objectives");
  }

  // T - Time-bound
  criteria.timeBound = validateTimeBound(goalData.deadline);
  if (!criteria.timeBound) {
    suggestions.push("Set a specific deadline to create urgency and focus");
  }

  // Calculate score
  const criteriaCount = Object.values(criteria).filter(Boolean).length;
  const score = (criteriaCount / 5) * 100;

  return {
    isValid: score >= 80, // Goal is valid if it meets 4+ criteria
    score,
    criteria,
    suggestions,
    warnings
  };
}

/**
 * Check if goal is Specific
 */
function validateSpecific(title: string, description: string): boolean {
  const text = `${title} ${description}`.toLowerCase();
  
  // Check for specific action words
  const actionWords = ['learn', 'complete', 'achieve', 'build', 'create', 'develop', 'improve', 'reduce', 'increase'];
  const hasActionWord = actionWords.some(word => text.includes(word));
  
  // Check for specific details (what, where, when, who, why)
  const hasDetails = text.length > 20 && description.length > 10;
  
  // Avoid vague words
  const vageWords = ['better', 'more', 'some', 'good', 'nice', 'things'];
  const hasVagueWords = vageWords.some(word => text.includes(word));
  
  return hasActionWord && hasDetails && !hasVagueWords;
}

/**
 * Check if goal is Measurable
 */
function validateMeasurable(title: string, description: string): boolean {
  const text = `${title} ${description}`.toLowerCase();
  
  // Check for numbers
  const hasNumbers = /\d+/.test(text);
  
  // Check for measurable units
  const units = ['kg', 'pounds', 'hours', 'days', 'weeks', 'months', 'times', 'percent', '%', 'dollars', '$'];
  const hasUnits = units.some(unit => text.includes(unit));
  
  // Check for measurable outcomes
  const measurableWords = ['complete', 'finish', 'reach', 'achieve', 'obtain', 'earn', 'save', 'lose', 'gain'];
  const hasMeasurableWords = measurableWords.some(word => text.includes(word));
  
  return hasNumbers || hasUnits || hasMeasurableWords;
}

/**
 * Check if goal is Achievable (basic heuristics)
 */
function validateAchievable(goalData: GoalData): boolean {
  const { title, description, deadline, priority } = goalData;
  const text = `${title} ${description}`.toLowerCase();
  
  // Check deadline reasonableness
  if (deadline) {
    const daysUntilDeadline = dayjs(deadline).diff(dayjs(), 'days');
    
    // Very short deadlines might be unrealistic for complex goals
    if (daysUntilDeadline < 1) return false;
    
    // Very long deadlines might lack urgency
    if (daysUntilDeadline > 365 * 2) return false;
  }
  
  // Check for unrealistic words
  const unrealisticWords = ['impossible', 'perfect', 'everything', 'never', 'always', 'instantly'];
  const hasUnrealisticWords = unrealisticWords.some(word => text.includes(word));
  
  // High priority goals should have reasonable scope
  if (priority === 'critical' && text.length < 30) {
    return false; // Critical goals should be well-defined
  }
  
  return !hasUnrealisticWords;
}

/**
 * Check if goal is Relevant
 */
function validateRelevant(goalData: GoalData): boolean {
  const { category, description } = goalData;
  
  // Basic relevance check based on category and description alignment
  const categoryKeywords: Record<string, string[]> = {
    health: ['fitness', 'diet', 'exercise', 'weight', 'wellness', 'nutrition', 'medical'],
    career: ['job', 'work', 'skill', 'promotion', 'salary', 'professional', 'leadership'],
    education: ['learn', 'study', 'course', 'degree', 'certification', 'knowledge', 'skill'],
    personal: ['habit', 'growth', 'development', 'mindfulness', 'hobby', 'relationship'],
    financial: ['money', 'save', 'invest', 'budget', 'income', 'debt', 'financial'],
    relationship: ['family', 'friend', 'social', 'communication', 'love', 'support']
  };
  
  const keywords = categoryKeywords[category] || [];
  const text = description.toLowerCase();
  
  // Check if description contains relevant keywords for the category
  const hasRelevantKeywords = keywords.some(keyword => text.includes(keyword));
  
  // Description should be substantial enough to show thought
  const hasSubstantialDescription = description.length > 20;
  
  return hasRelevantKeywords || hasSubstantialDescription;
}

/**
 * Check if goal is Time-bound
 */
function validateTimeBound(deadline?: string): boolean {
  if (!deadline) return false;
  
  const deadlineDate = dayjs(deadline);
  const now = dayjs();
  
  // Deadline should be in the future
  if (!deadlineDate.isAfter(now)) return false;
  
  // Deadline should be reasonable (not too far in the future)
  const daysUntilDeadline = deadlineDate.diff(now, 'days');
  if (daysUntilDeadline > 365 * 3) return false; // Max 3 years
  
  return true;
}

/**
 * Get SMART goal suggestions based on category
 */
export function getSMARTSuggestions(category: string): string[] {
  const suggestions: Record<string, string[]> = {
    health: [
      "Lose 10 pounds by exercising 30 minutes daily for 3 months",
      "Complete a 5K run in under 30 minutes by December 31st",
      "Eat 5 servings of vegetables daily for the next 2 months"
    ],
    career: [
      "Complete AWS certification course and pass exam within 6 months",
      "Lead 2 successful projects and get promoted to senior developer by year-end",
      "Increase sales by 20% through improved client relationships in Q4"
    ],
    education: [
      "Read 12 business books this year (1 per month)",
      "Complete online data science course with 90%+ grade in 4 months",
      "Learn Spanish to conversational level using Duolingo daily for 6 months"
    ],
    personal: [
      "Meditate for 10 minutes daily for the next 3 months",
      "Build a habit of journaling 5 minutes every morning for 2 months",
      "Declutter home by donating 100 items within 1 month"
    ],
    financial: [
      "Save $5,000 for emergency fund by setting aside $500 monthly",
      "Pay off $3,000 credit card debt in 8 months with $400 monthly payments",
      "Increase monthly income by $1,000 through freelance work within 6 months"
    ],
    relationship: [
      "Have one meaningful conversation with a family member weekly for 3 months",
      "Make 2 new professional connections monthly through networking events",
      "Plan and execute 1 fun activity with friends every 2 weeks for 6 months"
    ]
  };
  
  return suggestions[category] || [
    "Set a specific, measurable outcome with a clear deadline",
    "Break down large goals into smaller, actionable steps",
    "Ensure your goal aligns with your values and priorities"
  ];
}

/**
 * Analyze goal text and provide real-time feedback
 */
export function getGoalFeedback(title: string, description: string): {
  type: 'success' | 'warning' | 'info';
  message: string;
} {
  const text = `${title} ${description}`.toLowerCase();
  
  if (text.length < 10) {
    return {
      type: 'warning',
      message: 'Add more details to make your goal specific and clear'
    };
  }
  
  if (!/\d+/.test(text)) {
    return {
      type: 'info',
      message: 'Consider adding numbers or metrics to make your goal measurable'
    };
  }
  
  const vaguePhrases = ['get better', 'do more', 'be good', 'improve things'];
  if (vaguePhrases.some(phrase => text.includes(phrase))) {
    return {
      type: 'warning',
      message: 'Try to be more specific about what exactly you want to achieve'
    };
  }
  
  if (text.length > 50 && /\d+/.test(text)) {
    return {
      type: 'success',
      message: 'Great! Your goal looks specific and measurable'
    };
  }
  
  return {
    type: 'info',
    message: 'Keep adding details to make your goal SMART'
  };
}

export default {
  validateSMARTGoal,
  getSMARTSuggestions,
  getGoalFeedback
};